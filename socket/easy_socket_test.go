package socket

import (
	"testing"
	"time"
)

func TestRawServer(t *testing.T) {
	t.Log("TestRawServer")
	srv := &RawServer{Port: 6677}
	srv.Start()
}

func TestRawClient(t *testing.T) {
	t.Log("TestRawClient")
	TestManyRawClient("192.168.11.11", 6677, 100)
	time.Sleep(1800 * time.Second)
}

func TestWebServer(t *testing.T) {
	t.Log("TestWebServer")
	srv := &WebServer{Port: 6677}
	_, err := srv.Start(nil)
	if err != nil {
		t.Error(err)
	}
}

func TestWebClient(t *testing.T) {
	t.Log("TestWebClient")
}
