package socket

import (
	"fmt"
	"net"
	"time"
)

type RawClient struct {
	Host        string
	Port        int
	Name        string
	Monitor     *Monitor
	PrintDetail bool
}

func (cli *RawClient) Start() {
	if cli.PrintDetail {
		fmt.Printf("RawClient %s Start!\n", cli.Name)
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", cli.Host, cli.Port))
	if err != nil {
		fmt.Printf("[%s]创建连接失败，错误:%v\n", cli.Name, err)
		return
	}
	if cli.Monitor != nil {
		cli.Monitor.ValidNum <- 1
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("连接关闭失败，错误:%v\n\n", err)
		}
	}(conn)
	if cli.PrintDetail {
		fmt.Printf("[%s]连接服务端成功:%v\n", cli.Name, conn.RemoteAddr())
	}

	for i := 0; i < 99999999; i++ {
		nWrite, err := conn.Write([]byte(fmt.Sprintf("[%s]hello world(%d)", cli.Name, i)))
		if err != nil {
			fmt.Printf("发送消息失败:%v\n", err)
			return
		}
		if cli.Monitor != nil {
			cli.Monitor.BytesWrite <- nWrite
		}
		time.Sleep(200 * time.Millisecond)
		var buf [1024]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Printf("read failed, err:%v\n", err)
			return
		}
		if cli.Monitor != nil {
			cli.Monitor.BytesRead <- n
		}
		if cli.PrintDetail {
			if n > 0 {
				fmt.Printf("[%s] 收到服务端回复:%s\n", cli.Name, string(buf[:n]))
			} else {
				fmt.Printf("[%s] 没有收到服务端回复:\n", cli.Name)
			}
		}
	}
}
