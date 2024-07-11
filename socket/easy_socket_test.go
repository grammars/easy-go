package socket

import (
	"fmt"
	"testing"
)

func TestRawServer(t *testing.T) {
	fmt.Println("TestRawServer")
	srv := &SocketServer{Port: 6677}
	srv.Start()
}

func TestRawClient(t *testing.T) {
	fmt.Println("TestRawClient")
	cli := &SocketClient{Addr: "localhost", Port: 6677, Name: "好家伙"}
	cli.Start()
}
