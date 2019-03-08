package grpcurl

import (
	"fmt"
	"testing"
)

func TestEnerQuery(t *testing.T) {

	reply, err := EnerQuery("https://www.apkmirror.com/wp-content/themes/", "", -1, -1)
	fmt.Println("reply==", reply)
	if err != nil {
		t.Error(err)
	}
}

func TestQueryCount(t *testing.T) {
	QueryCount()
}
