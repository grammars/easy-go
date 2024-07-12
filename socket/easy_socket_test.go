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
	TestManyRawClient("localhost", 6677, 100)
	time.Sleep(600 * time.Second)
}
