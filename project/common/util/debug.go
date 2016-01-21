package util

func Assert(x interface{}) {
	if x == nil {
		return
	}

	b, ok := x.(bool)
	if ok == true {
		if b == true {
			return
		} else {
			panic(x)
		}
	}

	switch x.(type) {
	case uint8:
		Assert(0 == x.(uint8))
		return
	case int8:
		Assert(0 == x.(int8))
		return
	case uint16:
		Assert(0 == x.(uint16))
		return
	case int16:
		Assert(0 == x.(int16))
		return
	case uint32:
		Assert(0 == x.(uint32))
		return
	case int32:
		Assert(0 == x.(int32))
		return
	case uint64:
		Assert(0 == x.(uint64))
		return
	case int64:
		Assert(0 == x.(int64))
		return
	case uint:
		Assert(0 == x.(uint))
		return
	case int:
		Assert(0 == x.(int))
		return
	}

	panic(x)
}
