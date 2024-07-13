package socket

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net"
)

type RawServer struct {
	Port        int
	PrintDetail bool
	Monitor     *Monitor
}

func (srv *RawServer) Start() {
	slog.Info("开始启动SocketServer")
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", srv.Port))
	if err != nil {
		slog.Error("启动监听失败", "Error", err.Error())
		panic(err)
	}
	var conn net.Conn
	defer CloseConn(conn)
	for {
		conn, err := listen.Accept()
		if err != nil {
			slog.Error("接收连接失败", "Error", err)
			if srv.Monitor != nil {
				srv.Monitor.AcceptFailNum <- 1
			}
			continue
		}
		if srv.Monitor != nil {
			srv.Monitor.ValidNum <- 1
		}
		if srv.PrintDetail {
			slog.Info("有一个客户端连接我成功了", "RemoteAddr", conn.RemoteAddr())
		}
		go ReadWriteAsServer(conn, srv)
	}
}

func ReadWriteAsServer(conn net.Conn, srv *RawServer) {
	defer CloseConn(conn)
	reader := bufio.NewReader(conn)
	for {
		var buf [1024]byte
		n, err := reader.Read(buf[:])
		if err != nil && err != io.EOF {
			slog.Error("读取失败", "Error", err.Error())
			if srv.Monitor != nil {
				srv.Monitor.InvalidNum <- 1
			}
			break
		}
		if srv.Monitor != nil {
			srv.Monitor.BytesRead <- n
		}
		got := string(buf[:n])
		if srv.PrintDetail {
			slog.Info("接收到的数据", "数据", got)
		}
		nWrite, err := conn.Write([]byte("收到了：" + got))
		if err != nil {
			slog.Error("写给客户端失败", "Error", err.Error())
			return
		}
		if srv.Monitor != nil {
			srv.Monitor.BytesWrite <- nWrite
		}
		if srv.PrintDetail {
			slog.Info("回复客户端", "nWrite", nWrite)
		}
	}
}
