package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("Hello Gopher!")
	conn, err := net.Dial("tcp", "localhost:1888")
	if err != nil {
		fmt.Printf("创建连接失败，错误:%v\n", err)
		return
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("连接关闭失败，错误:%v\n\n", err)
		}
	}(conn)
	fmt.Printf("连接服务端成功:%v\n", conn.RemoteAddr())

	go func() {
		_, err := conn.Write([]byte("hello world"))
		if err != nil {
			fmt.Printf("发送消息失败:%v\n", err)
			return
		}
		time.Sleep(10 * time.Millisecond)
		var buf [1024]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Printf("read failed, err:%v\n", err)
			return
		}
		fmt.Println("收到服务端回复,", string(buf[:n]))
	}()

	time.Sleep(10 * time.Second)
}
