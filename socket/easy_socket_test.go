package socket

import (
	"fmt"
	"testing"
	"time"
)

func TestRawServer(t *testing.T) {
	fmt.Println("TestRawServer")
	srv := &RawServer{Port: 6677}
	srv.Start()
}

func TestRawClient(t *testing.T) {
	fmt.Println("TestRawClient")
	TestManyRawClient("192.168.11.11", 6677, 100)
	time.Sleep(1800 * time.Second)
}
