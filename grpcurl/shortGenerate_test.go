package grpcurl

import (
	"fmt"
	"testing"
)

func TestShortenerGenerate(t *testing.T) {
	reply, err := ShortGenerate("", "", "https://www.google.com")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("reply=", reply)
}
