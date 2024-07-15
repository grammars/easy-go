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

func TestWebServerTls(t *testing.T) {
	t.Log("TestWebServerTls")
	srv := &WebServer{Port: 0, PrintDetail: false, Monitor: CreateMonitorStart(),
		TLS: true, CrtFile: "E:\\gp\\assets\\server.crt", KeyFile: "E:\\gp\\assets\\server.key"}
	srv.StartDefault()
}

func TestWebClient(t *testing.T) {
	t.Log("TestWebClient")
	TestManyWebClient("192.168.11.11", 6677, false, 25)
	time.Sleep(1800 * time.Second)
}

func TestWebClientTls(t *testing.T) {
	t.Log("TestWebClientTls")
	TestManyWebClient("dev.ydwlgame.com", 0, true, 1)
	time.Sleep(1800 * time.Second)
}
