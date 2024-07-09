package main

import (
	"fmt"
	"net"
	"time"
)

func PlayRawSocket(index int) {
	fmt.Printf("PlayRawSocket %d begin!\n", index)
	conn, err := net.Dial("tcp", "localhost:1888")
	if err != nil {
		fmt.Printf("[%d]创建连接失败，错误:%v\n", index, err)
		return
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("连接关闭失败，错误:%v\n\n", err)
		}
	}(conn)
	fmt.Printf("%d连接服务端成功:%v\n", index, conn.RemoteAddr())

	for i := 0; i < 100; i++ {
		_, err := conn.Write([]byte(fmt.Sprintf("#%d#hello world(%d)", index, i)))
		if err != nil {
			fmt.Printf("发送消息失败:%v\n", err)
			return
		}
		time.Sleep(5000 * time.Millisecond)
		var buf [1024]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Printf("read failed, err:%v\n", err)
			return
		}
		if n > 0 {
			fmt.Printf("#%d# 收到服务端回复:%s\n", index, string(buf[:n]))
		} else {
			fmt.Printf("#%d# 没有收到服务端回复:\n")
		}
	}
}
