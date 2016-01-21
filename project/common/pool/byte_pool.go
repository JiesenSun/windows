package pool

import (
	"sync"
)

var (
	DEFAULT_BYTE_POOL_SIZE = 10000000
	IndexArray             = []int{8, 16, 32, 64, 128, 256, 1024, 2048}
	IndexNum               = len(IndexArray)
)

type ByteSlice struct {
	start int
	end   int
	index int
	Data  []byte
}

func bsCompare(v1, v2 interface{}) int {
	bs1 := v1.(*ByteSlice)
	bs2 := v2.(*ByteSlice)
	if bs1.end <= bs2.start {
		return -1
	} else if bs1.start >= bs2.end {
		return 1
	}
	return 0
}

type BytePool struct {
	*sync.Mutex
	data          []byte
	start         int
	end           int
	size          int
	usedList      []*StaticList
	freeList      []*StaticList
	byteSlicePool *Pools
}

func NewBytePool() *BytePool {
	bytePool := &BytePool{
		Mutex:         &sync.Mutex{},
		data:          make([]byte, DEFAULT_BYTE_POOL_SIZE),
		size:          DEFAULT_BYTE_POOL_SIZE,
		end:           DEFAULT_BYTE_POOL_SIZE,
		usedList:      []*StaticList{},
		freeList:      []*StaticList{},
		byteSlicePool: NewPools(BytePool{}, DEFAULT_OBJECT_SIZE),
	}
	for i := 0; i < IndexNum; i++ {
		bytePool.usedList[i] = NewStaticList()
		bytePool.freeList[i] = NewStaticList()
	}
	return bytePool
}

func (this *BytePool) Index(size int) int {
	for i := 0; i < IndexNum; i++ {
		if size <= IndexArray[i] {
			return i
		}
	}
	return IndexNum
}

func (this *BytePool) Get(size int) *ByteSlice {
	index := this.Index(size)

	// 申请大块内存 size > 2048  直接系统分配
	if index == IndexNum {
		bs := this.byteSlicePool.Get().(*ByteSlice)
		bs.start = 0
		bs.end = size
		bs.Data = make([]byte, size)
		return bs
	}
	// 小块内存 从链块中取
	bs := this.freeList[index].PopFront()
	if bs == nil && index == IndexNum-1 { // 链块为空
		bs := this.byteSlicePool.Get().(*ByteSlice)
		this.Lock()
		bs.start = this.start
		bs.end = bs.start + size
		if bs.end > this.end {
			this.byteSlicePool.Put(bs)
			this.Unlock()
			return nil
		}
		this.start = bs.end
		this.Unlock()
		bs.Data = this.data[bs.start:bs.end]
		return bs
	} else if bs == nil {
		bs := this.Get(IndexArray[index+1])
		if bs == nil {
			return nil
		}
		smallBS := this.byteSlicePool.Get().(*ByteSlice)
		smallBS.start = bs.start
		smallBS.end = smallBS.start + IndexArray[index]
		bs.start = smallBS.end
		this.freeList[index].PushFront(bs)
		return smallBS
	} else {
		return bs.(*ByteSlice)
	}
	return nil
}

func (this *BytePool) Put(bs *ByteSlice) {
	index := this.Index(bs.end - bs.start)

	if index == IndexNum {
		bs.Data = nil
		this.byteSlicePool.Put(bs)
	} else {
		this.freeList[index].PushFront(bs)
	}
}

type BytePools struct {
	pool []*BytePool
}

func NewBytePools() *BytePools {
	return &BytePools{
		pool: []*BytePool{},
	}
}

func (this *BytePools) Get(size int) *ByteSlice {
	var bs *ByteSlice = nil

	for i, pool := range this.pool {
		bs = pool.Get(size)
		if bs != nil {
			return bs
		}
		bs.index = i
	}
	this.pool = append(this.pool, NewBytePool())
	return this.Get(size)
}

func (this *BytePools) Put(bs *ByteSlice) {
	this.pool[bs.index].Put(bs)
}
