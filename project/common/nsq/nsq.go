package nsq

import (
	"strings"
	"sync"

	"github.com/bitly/go-nsq"

	"sirendaou.com/duserver/common/errors"
	sync_ "sirendaou.com/duserver/common/sync"
	"sirendaou.com/duserver/common/syslog"
)

const (
	MSG_CHAN_SIZE           = 10000
	PRODUCER_COUNT_PER_ADDR = 1
	CONSUMER_COUNT_PER_ADDR = 3
)

var (
	g_nsqdMgr *nsqdMgr = &nsqdMgr{
		rwMutex:         &sync.RWMutex{},
		topicProducer:   make(map[string]*Producer),
		defaultProducer: nil,
		consumer:        nil,
	}
)

type Producer struct {
	addrs      string
	producerCh chan *nsq.Producer
}

func NewProducer(addrs string) *Producer {
	config := nsq.NewConfig()
	config.DefaultRequeueDelay = 0

	addrSlice := strings.Split(addrs, ",")
	produceCount := len(addrSlice) * PRODUCER_COUNT_PER_ADDR
	producerCh := make(chan *nsq.Producer, produceCount)

	for _, addr := range addrSlice {
		for i := 0; i < PRODUCER_COUNT_PER_ADDR; i++ {
			producer, err := nsq.NewProducer(addr, config)
			if err != nil {
				panic(err)
			}
			producerCh <- producer
		}
	}

	return &Producer{
		addrs:      addrs,
		producerCh: producerCh,
	}
}

func (this *Producer) Publish(topic string, body []byte) error {
	producer := <-this.producerCh
	defer func() {
		this.producerCh <- producer
	}()
	return producer.Publish(topic, body)
}

func (this *Producer) Close() {
	count := cap(this.producerCh)
	for i := 0; i < count; i++ {
		p := <-this.producerCh
		p.Stop()
	}
}

func Publish(topic string, body []byte) error {
	g_nsqdMgr.rwMutex.RLock()
	p, ok := g_nsqdMgr.topicProducer[topic]
	g_nsqdMgr.rwMutex.RUnlock()
	if ok {
		return p.Publish(topic, body)
	}
	return g_nsqdMgr.defaultProducer.Publish(topic, body)
}

type nsqdMgr struct {
	addrs           string
	rwMutex         *sync.RWMutex
	topicProducer   map[string]*Producer
	defaultProducer *Producer
	consumer        *ConsumerT
}

func Init(addrs string) {
	g_nsqdMgr.addrs = addrs
	g_nsqdMgr.defaultProducer = NewProducer(addrs)
}

// message适配代码
type Message nsq.Message
type Handler interface {
	HandleMessage(mesage *Message) error
}
type ConsumerT struct {
	consumers []*nsq.Consumer
	handle    Handler
	msgCh     chan *Message
	waitGroup *sync_.WaitGroup
}

func (this *ConsumerT) HandleMessage(message *nsq.Message) error {
	this.msgCh <- (*Message)(message)
	return nil
}

func Consumer(topic, channel string, handle Handler) (*ConsumerT, error) {
	return ConsumerGO(topic, channel, 1, handle)
}

func ConsumerGO(topic, channel string, goCount uint, handle Handler) (*ConsumerT, error) {
	msgHandle := &ConsumerT{
		consumers: []*nsq.Consumer{},
		handle:    handle,
		msgCh:     make(chan *Message, MSG_CHAN_SIZE),
		waitGroup: sync_.NewWaitGroup(),
	}
	addrSlice := strings.Split(g_nsqdMgr.addrs, ",")
	for _, addr := range addrSlice {
		for i := 0; i < CONSUMER_COUNT_PER_ADDR; i++ {
			consumer, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())
			if err != nil {
				return nil, errors.As(err, topic, channel)
			}
			consumer.SetLogger(nil, nsq.LogLevelInfo)
			consumer.AddHandler(msgHandle)
			if err := consumer.ConnectToNSQD(addr); err != nil {
				return nil, errors.As(err, topic, channel, g_nsqdMgr.addrs)
			}
			msgHandle.consumers = append(msgHandle.consumers, consumer)
		}
	}
	g_nsqdMgr.consumer = msgHandle
	for i := 0; i < int(goCount); i++ {
		go msgHandle.work()
	}
	return msgHandle, nil
}

func (this *ConsumerT) work() {
	this.waitGroup.Add(1)
	defer this.waitGroup.Done()
	exitNotify := this.waitGroup.ExitNotify()
	for {
		select {
		case <-exitNotify:
			return
		case msg := <-this.msgCh:
			if err := this.handle.HandleMessage(msg); err != nil {
				syslog.Info(err, msg)
			}
		}
	}
}

func Deinit() {
	g_nsqdMgr.defaultProducer.Close()
	g_nsqdMgr.rwMutex.RLock()
	for _, p := range g_nsqdMgr.topicProducer {
		p.Close()
	}
	g_nsqdMgr.rwMutex.RUnlock()
	if g_nsqdMgr.consumer != nil {
		for _, c := range g_nsqdMgr.consumer.consumers {
			c.Stop()
		}
		g_nsqdMgr.consumer.waitGroup.Wait()
	}
}
