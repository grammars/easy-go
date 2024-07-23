package socket

import (
	"fmt"
	socket "github.com/grammars/easy-go/socket/codec"
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
	// reader := bufio.NewReader(conn)
	for {
		cr, err := srv.Decoder.Decode(conn)
		if err != nil {
			slog.Error("读取失败", "Error", err.Error())
			break
		}
		slog.Info("本帧长度", cr.FrameLength, "HeaderBytes", cr.HeaderBytes, "HeaderBytes", cr.HeaderBytes)
		//var buf [1024]byte
		//n, err := reader.Read(buf[:])
		//if err != nil && err != io.EOF {
		//	slog.Error("读取失败", "Error", err.Error())
		//	break
		//}
		if srv.Monitor != nil {
			srv.Monitor.BytesRead <- cr.FrameLength
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
