package socket

import (
	"fmt"
	"time"
)

func TestManyRawClient(host string, port int, clientNum int) []*RawClient {
	monitor := CreateMonitorStart()
	var clients []*RawClient
	for i := 0; i < clientNum; i++ {
		cli := &RawClient{Host: host, Port: port, Name: fmt.Sprintf("小家伙%d", i),
			Monitor: monitor, PrintDetail: clientNum == 1}
		go cli.Start()
		clients = append(clients, cli)
		time.Sleep(20 * time.Millisecond)
	}
	return clients
}

func TestManyWebClient(host string, port int, clientNum int) []*WebClient {
	monitor := CreateMonitorStart()
	var clients []*WebClient
	for i := 0; i < clientNum; i++ {
		cli := &WebClient{Host: host, Port: port, Name: fmt.Sprintf("好家伙%d", i),
			Monitor: monitor, PrintDetail: clientNum == 1}
		go cli.Start()
		clients = append(clients, cli)
		time.Sleep(20 * time.Millisecond)
	}
	return clients
}
