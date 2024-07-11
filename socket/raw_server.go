package socket

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

type RawServer struct {
	Port int
}

func (srv *RawServer) Start() {
	fmt.Println("开始启动SocketServer")

	monitor := Monitor{}
	go monitor.Start()

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", srv.Port))
	if err != nil {
		fmt.Printf("启动监听失败，错误:%v\n\n", err)
		panic(err)
	}
	var conn net.Conn
	defer CloseConn(conn)
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("接收连接失败:%v\n", err)
			continue
		}
		fmt.Printf("有一个客户端连接我成功了，来自:%v\n", conn.RemoteAddr())
		go ReadWriteAsServer(conn, &monitor)
	}
}

func ReadWriteAsServer(conn net.Conn, monitor *Monitor) {
	defer CloseConn(conn)
	reader := bufio.NewReader(conn)
	for {
		var buf [1024]byte
		n, err := reader.Read(buf[:])
		if err != nil && err != io.EOF {
			//fmt.Printf("读取失败,err:%v\n", err)
			monitor.InvalidNum <- 1
			break
		}
		monitor.BytesRead <- n
		got := string(buf[:n])
		fmt.Println("接收到的数据: ", got)
		nWrite, err := conn.Write([]byte("收到了：" + got))
		if err != nil {
			return
		}
		fmt.Println("回复:nWrite=", nWrite)
	}
}
