package ego

import "github.com/grammars/easy-go/socket"

func Return[T any](boolExpression bool, trueReturnValue, falseReturnValue T) T {
	if boolExpression {
		return trueReturnValue
	} else {
		return falseReturnValue
	}
}

func ReturnByFunc[T any](boolExpression bool, trueFuncForReturnValue, falseFuncForReturnValue func() T) T {
	if boolExpression {
		return trueFuncForReturnValue()
	} else {
		return falseFuncForReturnValue()
	}
}

func CreateSocketServer(port int) *socket.SocketServer {
	return &socket.SocketServer{Port: port}
}

func CreateSocketClient(addr string, port int, name string) *socket.SocketClient {
	return &socket.SocketClient{Addr: addr, Port: port, Name: name}
}
