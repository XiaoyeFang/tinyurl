package grpcurl

import (
	"fmt"
	"testing"
)

func TestEnerDelete(t *testing.T) {
	reply, err := EnerDelete("http://www.baidu.com", "http://192.168.9.125:8080/ZdIVR9")
	fmt.Println(reply)
	if err != nil {
		t.Error(err)
	}
}
