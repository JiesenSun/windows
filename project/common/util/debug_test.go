package util

import (
	"fmt"
	"testing"
)

func TestAssert(t *testing.T) {
	_, err := fmt.Println("hello world")
	Assert(err)
	Assert(1 == 1)
	Assert(0)
}
