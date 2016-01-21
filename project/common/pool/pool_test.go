package pool

import (
	"fmt"
	"testing"
)

type T struct {
	a int
}

func TestPoolMgr(t *testing.T) {
	typName := "T"
	RegisterType(typName, T{}, 100)

	obj := Get(typName).(*T)
	obj.a = 1
	fmt.Println(obj)
}
