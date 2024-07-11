package ego

import "testing"

func TestAny(t *testing.T) {
	Practice.Echo()
}

func TestSocketServer(t *testing.T) {
	srv := Socket.CreateRawServer(2345)
	srv.Start()
}

func TestSocketClient(t *testing.T) {
	srv := Socket.CreateRawClient("localhost", 2345, "单元测试SCli")
	srv.Start()
}
