package protocol

import (
	"encoding/json"
	"fmt"
	"testing"
)

type TestObject struct{}

func TestRegisterObject(t *testing.T) {
	RegisterObject(1, TestObject{})

	obj, ok := GetObject(1).(*TestObject)
	fmt.Println(ok)

	if err := json.Unmarshal([]byte("{}"), obj); err != nil {
		fmt.Println(err)
	}
}
