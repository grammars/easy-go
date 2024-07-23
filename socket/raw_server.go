package socket

import (
	"fmt"
	socket "github.com/grammars/easy-go/socket/codec"
	"io"
	"log/slog"
	"net"
	"time"
)

type RawServer[VD any] struct {
	Port        int
	PrintDetail bool

	StartTime  time.Time
	Monitor    *Monitor
	Handler    VisitorServerHandler[VD]
	VisitorMap *VisitorMap[VD]
	Decoder    socket.FrameDecoder
}

func (srv *RawServer[VD]) Start() {
	slog.Info("开始启动SocketServer")
	srv.StartTime = time.Now()
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", srv.Port))
	if err != nil {
		slog.Error("启动监听失败", "Error", err.Error())
		panic(err)
	}
	srv.VisitorMap = CreateVisitorMap[VD](srv)
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
		if srv.PrintDetail {
			slog.Info("有一个客户端连接我成功了", "RemoteAddr", conn.RemoteAddr())
		}
		go ReadWriteAsServer(conn, srv)
	}
}

func (srv *RawServer[VD]) GetStartTime() time.Time {
	return srv.StartTime
}

func (srv *RawServer[VD]) appendVisitor(conn net.Conn) *Visitor[VD] {
	visitor := srv.VisitorMap.Append(conn)
	slog.Info("Accept客户端", "uid", visitor.uid, "index", visitor.index, "addr", conn.RemoteAddr())
	if srv.PrintDetail {
		srv.VisitorMap.Print()
	}
	if srv.Monitor != nil {
		srv.Monitor.ValidNum <- 1
	}
	return visitor
}

func (srv *RawServer[VD]) removeVisitor(visitorUid uint64) {
	srv.VisitorMap.Remove(visitorUid)
	slog.Info("Remove客户端", "visitorUid", visitorUid)
	if srv.PrintDetail {
		srv.VisitorMap.Print()
	}
	srv.Monitor.InvalidNum <- 1
}

func ReadWriteAsServer[VD any](conn net.Conn, srv *RawServer[VD]) {
	visitor := srv.appendVisitor(conn)
	if srv.Handler != nil {
		srv.Handler.OnConnect(visitor)
	}
	defer srv.removeVisitor(visitor.uid)
	defer func() {
		if srv.Handler != nil {
			srv.Handler.OnDisconnect(visitor)
		}
	}()
	defer CloseConn(conn)
	for {
		cr, err := srv.Decoder.Decode(conn)
		if err != nil {
			slog.Error("读取失败", "Error", err.Error())
			break
		}
		slog.Info("本帧长度", "FrameLength", cr.FrameLength)

		if cr.Overflow {
			slog.Info("即将执行数据帧溢出后处理")
			err := conn.SetDeadline(time.Now().Add(time.Second * 1))
			if err != nil {
				break
			}
			leftAll, err := io.ReadAll(conn)
			if err != nil {
				break
			}
			slog.Info("读完剩余的", "leftAll长度", len(leftAll))
		}

		if srv.Monitor != nil {
			srv.Monitor.BytesRead <- cr.FrameLength
		}

		if srv.Handler != nil {
			srv.Handler.OnMessage(visitor, cr)
		}
		//got := string(buf[:n])
		//if srv.PrintDetail {
		//	slog.Info("接收到的数据", "数据", got)
		//}
		//nWrite, err := conn.Write([]byte("收到了：" + got))
		//if err != nil {
		//	slog.Error("写给客户端失败", "Error", err.Error())
		//	return
		//}
		//if srv.Monitor != nil {
		//	srv.Monitor.BytesWrite <- nWrite
		//}
		//if srv.PrintDetail {
		//	slog.Info("回复客户端", "nWrite", nWrite)
		//}
	}
}
