package socket

func CreateRawServer(port int) *RawServer {
	return &RawServer{Port: port}
}

func CreateRawClient(addr string, port int, name string) *RawClient {
	return &RawClient{Addr: addr, Port: port, Name: name}
}
