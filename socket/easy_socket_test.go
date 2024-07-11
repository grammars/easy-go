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
	monitor := &Monitor{}
	go monitor.Start()
	for i := 0; i < 3; i++ {
		cli := &RawClient{Addr: "localhost", Port: 6677, Name: fmt.Sprintf("好家伙%d", i), Monitor: monitor}
		go cli.Start()
	}
	time.Sleep(20 * time.Second)
}
