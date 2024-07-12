package socket

import (
	"fmt"
	"time"
)

func CreateRawServer(port int) *RawServer {
	return &RawServer{Port: port}
}

func CreateRawClient(host string, port int, name string) *RawClient {
	return &RawClient{Host: host, Port: port, Name: name}
}

func TestManyRawClient(host string, port int, clientNum int) []*RawClient {
	monitor := &Monitor{}
	go monitor.Start()
	var clients []*RawClient
	for i := 0; i < clientNum; i++ {
		cli := &RawClient{Host: host, Port: port, Name: fmt.Sprintf("好家伙%d", i), Monitor: monitor}
		go cli.Start()
		clients = append(clients, cli)
		time.Sleep(20 * time.Millisecond)
	}
	return clients
}
