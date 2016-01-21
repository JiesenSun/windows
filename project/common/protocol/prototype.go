package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"sync"
)

/*
* program const
 */
const (
	DATETIME_FMT = "2006-01-02 15:04:05"
	BUFFER_SIZE  = 2048
)

/*
*  command const
 */
const (
	CLIENT_CMD_HEARTBEAT           = 10000
	CLIENT_CMD_USER_REGISTER       = 10001
	CLIENT_CMD_USER_LOGIN          = 10002
	CLIENT_CMD_USER_RESET_PWD      = 10003
	CLIENT_CMD_USER_LOGIN_CONFLICT = 10004

	CLIENT_CMD_SEND_USER_MSG = 10005
	CLIENT_CMD_RECV_USER_MSG = 10006
)

/*
*   与客户端通信协议
*	包格式:包头+包体
 */
const (
	DATA_PACKAGE_MAX_SIZE  = 2048
	DATA_PACKAGE_MIN_SIZE  = 12
	DATA_PACKAGE_HEAD_SIZE = 12
	DATA_PACKAGE_TAIL_SIZE = 20
)

/*
* memory manager
 */

var (
	g_datapkg_pool = &sync.Pool{New: func() interface{} { return &DataPackage{} }}
	g_buffer_pool  = &sync.Pool{New: func() interface{} { return NewBuffer() }}

	g_object_map = make(map[int]reflect.Type)
)

func RegisterObject(cmd int, object interface{}) {
	obj := reflect.Indirect(reflect.ValueOf(object))
	g_object_map[cmd] = obj.Type()
}

func GetObject(cmd int) interface{} {
	typ, ok := g_object_map[cmd]
	if ok == false {
		return nil
	}

	ind := reflect.New(typ)
	return ind.Interface()
}

func GetDataPackage() *DataPackage {
	return g_datapkg_pool.Get().(*DataPackage)
}

func PutDataPackage(dp *DataPackage) {
	g_datapkg_pool.Put(dp)
}

func GetBuffer() *Buffer {
	buf := g_buffer_pool.Get().(*Buffer)
	buf.Reset()
	return buf
}

func PutBuffer(buf *Buffer) {
	g_buffer_pool.Put(buf)
}

type Buffer struct {
	*bytes.Buffer
	buf [BUFFER_SIZE]byte
}

func NewBuffer() *Buffer {
	buf := &Buffer{}
	buf.Buffer = bytes.NewBuffer(buf.buf[:])
	buf.Reset()
	return buf
}

func (buf *Buffer) Buff() []byte {
	return buf.buf[:]
}

type PackageHead struct {
	PkgLen     uint16 // head + body
	Version    uint16
	Command    uint32
	SequenceID uint32
}

type PackageTail struct {
	UID  uint64
	SID  uint32
	IP   uint32
	Port uint32
}

type DataPackage struct {
	Head PackageHead
	Body [DATA_PACKAGE_MAX_SIZE]byte
	Tail PackageTail
	Data interface{}
}

func (this *DataPackage) Unpackage(data []byte) error {
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &this.Head); err != nil {
		return err
	}
	dataLen := len(data)

	if dataLen == DATA_PACKAGE_HEAD_SIZE {
		return nil
	}

	if (this.Head.PkgLen != uint16(dataLen)) && (this.Head.PkgLen+DATA_PACKAGE_TAIL_SIZE != uint16(dataLen)) {
		return errors.New("data package len errors:")
	}

	if this.Head.PkgLen > DATA_PACKAGE_MAX_SIZE+DATA_PACKAGE_HEAD_SIZE {
		return errors.New("package length too large")
	}
	tempBuf := this.Body[:this.Head.PkgLen-DATA_PACKAGE_HEAD_SIZE]
	if err := binary.Read(buf, binary.BigEndian, &tempBuf); err != nil {
		return err
	}

	if this.Head.PkgLen+DATA_PACKAGE_TAIL_SIZE == uint16(dataLen) {
		if err := binary.Read(buf, binary.BigEndian, &this.Tail); err != nil {
			return err
		}
	}

	return nil
}

func (this *DataPackage) Package(w io.Writer) error {
	return Package(w, &this.Head, this.BodyData(), &this.Tail)
}

func (this *DataPackage) BodyData() []byte {
	return this.Body[:this.Head.PkgLen-DATA_PACKAGE_HEAD_SIZE]
}

func Package(w io.Writer, head *PackageHead, body []byte, tail *PackageTail) error {
	if err := binary.Write(w, binary.BigEndian, *head); err != nil {
		return err
	}
	if body != nil && len(body) != 0 {
		if err := binary.Write(w, binary.BigEndian, body); err != nil {
			return err
		}
	}
	if err := binary.Write(w, binary.BigEndian, *tail); err != nil {
		return err
	}

	return nil
}

func JsonMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
func JsonUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

//===================udp prototype=====================
type UdpPackage struct {
	SID        uint32
	BodyLength uint32
	Body       [0]byte
}
