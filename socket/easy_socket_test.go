package socket

import (
	"fmt"
	"testing"
)

func TestRawServer(t *testing.T) {
	fmt.Println("TestRawServer")
	srv := &RawServer{Port: 6677}
	srv.Start()
}

func TestRawClient(t *testing.T) {
	fmt.Println("TestRawClient")
	cli := &RawClient{Addr: "localhost", Port: 6677, Name: "好家伙"}
	cli.Start()
}
