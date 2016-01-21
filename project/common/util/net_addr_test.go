package util

import (
	"fmt"
	"testing"
)

func TestIPConv(t *testing.T) {
	ipInt := IPToInt("192.168.20.25")
	fmt.Println(ipInt)
	fmt.Println(IPToStr(ipInt))
	ipInt = IPToInt("192.168.20.25:9000")
	fmt.Println(ipInt)
	fmt.Println(IPToStr(ipInt))
}
