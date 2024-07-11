package ego

import "github.com/grammars/easy-go/socket"

func CreateSocketServer(port int) *socket.SocketServer {
	return &socket.SocketServer{Port: port}
}

func CreateSocketClient(addr string, port int, name string) *socket.SocketClient {
	return &socket.SocketClient{Addr: addr, Port: port, Name: name}
}
