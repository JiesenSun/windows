package main

import (
	_ "project/common/syscall"
	"project/gateway"
)

func main() {
	gateway.StartServer()
}
