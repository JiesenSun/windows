package pool

import (
	"reflect"
	"sync"
)

var (
	DEFAULT_OBJECT_SIZE = 10000
)

type Pool struct {
	mutex     *sync.Mutex
	rwMutex   *sync.RWMutex
	typ       reflect.Type
	data      reflect.Value
	addrIndex map[interface{}]int
	next      []int
	head      int
	size      int
}

func NewPool(obj interface{}, size int) *Pool {
	ind := reflect.ValueOf(obj)
	typ := reflect.Indirect(ind).Type()
	if typ.Kind() == reflect.Chan || typ.Kind() == reflect.Map || typ.Kind() == reflect.Slice {
		println("not support chan and map slice!!!")
		return nil
	}

	pool := &Pool{
		mutex:     &sync.Mutex{},
		rwMutex:   &sync.RWMutex{},
		typ:       typ,
		data:      reflect.MakeSlice(reflect.SliceOf(typ), size, size),
		addrIndex: make(map[interface{}]int),
		next:      make([]int, size, size),
		head:      0,
		size:      size,
	}

	for i := 0; i < size; i++ {
		pool.next[i] = i + 1
	}
	return pool
}

func (this *Pool) Full() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.head == this.size
}

func (this *Pool) Get() interface{} {
	// 取出一个节点地址
	this.mutex.Lock()
	if this.head == this.size {
		// 申请内存已耗尽 动态扩展TODO
		this.mutex.Unlock()
		//return reflect.New(this.typ).Interface()
		return nil
	}
	free := this.head
	nextFree := this.next[free]
	this.head = nextFree
	this.mutex.Unlock()

	// 保存节点地址
	obj := this.data.Index(free).Addr().Interface()
	this.rwMutex.Lock()
	this.addrIndex[obj] = free
	this.rwMutex.Unlock()
	return obj
}

func (this *Pool) Put(obj interface{}) {
	this.rwMutex.Lock()
	index, ok := this.addrIndex[obj]
	delete(this.addrIndex, obj)
	this.rwMutex.Unlock()
	if ok == false {
		return
	}
	this.mutex.Lock()
	this.next[index] = this.head
	this.head = index
	this.mutex.Unlock()

	return
}

type Pools struct {
	obj     interface{}
	size    int
	pools   []*Pool
	rwMutex *sync.RWMutex
}

func NewPools(obj interface{}, size int) *Pools {
	return &Pools{
		obj:     obj,
		size:    size,
		pools:   []*Pool{},
		rwMutex: &sync.RWMutex{},
	}
}
func (this *Pools) Get() interface{} {
	var obj interface{} = nil

	this.rwMutex.RLock()
	for _, pool := range this.pools {
		if pool.Full() {
			continue
		}
		obj = pool.Get()
		break
	}
	this.rwMutex.RUnlock()

	if obj == nil {
		this.rwMutex.Lock()
		this.pools = append(this.pools, NewPool(this.obj, this.size))
		this.rwMutex.Unlock()
		obj = this.Get()
	}
	return obj
}

func (this *Pools) Put(obj interface{}) {
	this.rwMutex.RLock()
	for _, pool := range this.pools {
		pool.Put(obj)
	}
	this.rwMutex.RUnlock()
}

type PoolsMgr struct {
	rwMutex *sync.RWMutex
	pool    map[string]*Pools
}

var (
	g_default_pool = &PoolsMgr{
		rwMutex: &sync.RWMutex{},
		pool:    make(map[string]*Pools),
	}
)

func RegisterType(typ string, obj interface{}, size int) {
	g_default_pool.rwMutex.Lock()
	defer g_default_pool.rwMutex.Unlock()
	_, ok := g_default_pool.pool[typ]
	if ok {
		return
	}

	g_default_pool.pool[typ] = NewPools(obj, size)
}

func Get(typ string) interface{} {
	g_default_pool.rwMutex.RLock()
	pool, ok := g_default_pool.pool[typ]
	if !ok {
		g_default_pool.rwMutex.RUnlock()
		return nil
	}
	g_default_pool.rwMutex.RUnlock()

	return pool.Get()
}

func Put(typ string, obj interface{}) {
	g_default_pool.rwMutex.RLock()
	pool, ok := g_default_pool.pool[typ]
	if !ok {
		g_default_pool.rwMutex.RUnlock()
		return
	}
	g_default_pool.rwMutex.RUnlock()

	pool.Put(obj)
}
