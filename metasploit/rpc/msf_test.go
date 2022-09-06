package rpc

import (
	"testing"
)

func TestMsf(t *testing.T) {
	a, _ := New("192.168.1.128:55552", "123", "2222")
	a.send(loginReq{}, &loginRes{})
}
