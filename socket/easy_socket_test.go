package socket

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
	"time"
)

func TestRawServer(t *testing.T) {
	t.Log("TestRawServer")
	srv := &RawServer{Port: 6677, Monitor: CreateMonitorStart("测试")}
	srv.Start()
}

func TestRawClient(t *testing.T) {
	t.Log("TestRawClient")
	TestManyRawClient("192.168.11.11", 6677, 100)
	time.Sleep(1800 * time.Second)
}

func TestWebServer(t *testing.T) {
	t.Log("TestWebServer")
	srv := &WebServer[any]{Port: 6677, PrintDetail: false, Monitor: CreateMonitorStart("测试")}
	srv.StartDefault()
}

func TestWebServerTls(t *testing.T) {
	t.Log("TestWebServerTls")
	srv := &WebServer[any]{Port: 0, PrintDetail: false, Monitor: CreateMonitorStart("测试"),
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

func TestEndian(t *testing.T) {
	t.Log("TestEndian")
	var num int32 = 305419896 // 0x12345678

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		t.Log("binary.Write failed:", err)
		return
	}
	t.Log("BigEndian -> ", fmt.Sprintf("0x%x", buf.Bytes()))

	buf = new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, num)
	if err != nil {
		t.Log("binary.Write failed:", err)
		return
	}
	t.Log("LittleEndian -> ", fmt.Sprintf("0x%x", buf.Bytes()))
}
