package socket

import (
	"testing"
	"time"
)

func TestRawServer(t *testing.T) {
	t.Log("TestRawServer")
	srv := &RawServer{Port: 6677, Monitor: CreateMonitorStart()}
	srv.Start()
}

func TestRawClient(t *testing.T) {
	t.Log("TestRawClient")
	TestManyRawClient("192.168.11.11", 6677, 100)
	time.Sleep(1800 * time.Second)
}

func TestWebServer(t *testing.T) {
	t.Log("TestWebServer")
	srv := &WebServer{Port: 6677, PrintDetail: false, Monitor: CreateMonitorStart()}
	srv.StartDefault()
}

func TestWebClient(t *testing.T) {
	t.Log("TestWebClient")
	TestManyWebClient("192.168.11.11", 6677, 25)
	time.Sleep(1800 * time.Second)
}
