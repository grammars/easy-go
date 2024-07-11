package ego

import "testing"

func TestSocketServer(t *testing.T) {
	srv := CreateSocketServer(2345)
	srv.Start()
}

func TestSocketClient(t *testing.T) {
	srv := CreateSocketClient("localhost", 2345, "单元测试SCli")
	srv.Start()
}
