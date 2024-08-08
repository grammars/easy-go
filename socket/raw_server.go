package socket

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"
)

type RawServer[VD any] struct {
	Port int

	StartTime         time.Time
	Monitor           *Monitor
	Handler           VisitorServerHandler[VD]
	VisitorMap        *VisitorMap[VD]
	Decoder           FrameDecoder[VD]
	FrameBrokenDumpMs time.Duration // 当数据帧损坏时，倾倒时间(ms) 小于等于0则采取断开策略
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
		if LogLevel <= 1 {
			slog.Info("raw_server.go 客户端连接成功", "RemoteAddr", conn.RemoteAddr())
		}
		go ReadWriteAsServer(conn, srv)
	}
}

func (srv *RawServer[VD]) GetStartTime() time.Time {
	return srv.StartTime
}

type RawVisitorConnection[VD any] struct {
	conn       net.Conn
	WriteMutex sync.Mutex
	srv        *RawServer[VD]
}

func (rvc *RawVisitorConnection[VD]) RemoteAddr() net.Addr {
	return rvc.conn.RemoteAddr()
}

func (rvc *RawVisitorConnection[VD]) Write(b []byte) (int, error) {
	n, err := rvc.conn.Write(b)
	if rvc.srv.Monitor != nil {
		rvc.srv.Monitor.BytesWrite <- n
	}
	return n, err
}

func (rvc *RawVisitorConnection[VD]) WriteSafe(b []byte) (int, error) {
	rvc.WriteMutex.Lock()
	defer rvc.WriteMutex.Unlock()
	return rvc.Write(b)
}

func (srv *RawServer[VD]) appendVisitor(conn net.Conn) *Visitor[VD] {
	rvc := &RawVisitorConnection[VD]{conn: conn, srv: srv}
	visitor := srv.VisitorMap.Append(rvc)
	if LogLevel <= 0 {
		slog.Info("raw_server.go Accept访问者", "uid", visitor.Uid, "index", visitor.index, "addr", conn.RemoteAddr())
		srv.VisitorMap.Print()
	}
	if srv.Monitor != nil {
		srv.Monitor.ValidNum <- 1
	}
	return visitor
}

func (srv *RawServer[VD]) removeVisitor(visitorUid uint64) {
	srv.VisitorMap.Remove(visitorUid)
	if LogLevel <= 0 {
		slog.Info("raw_server.go Remove访问者", "visitorUid", visitorUid)
		srv.VisitorMap.Print()
	}
	srv.Monitor.InvalidNum <- 1
}

func ReadWriteAsServer[VD any](conn net.Conn, srv *RawServer[VD]) {
	visitor := srv.appendVisitor(conn)
	if srv.Handler != nil {
		srv.Handler.OnConnect(visitor)
	}
	defer srv.removeVisitor(visitor.Uid)
	defer func() {
		if srv.Handler != nil {
			srv.Handler.OnDisconnect(visitor)
		}
	}()
	defer CloseConn(conn)
	for {
		cr, err := srv.Decoder.Decode(visitor, conn)
		if err != nil {
			slog.Error("读取失败", "Error", err.Error())
			break
		}
		if LogLevel <= 0 {
			slog.Info("本帧长度", "FrameLength", cr.FrameLength)
		}

		if cr.Overflow {
			brokenAction := "即将执行数据帧溢出后处理"
			if srv.FrameBrokenDumpMs <= 0 {
				slog.Info(brokenAction, "strategy", "FrameBrokenDumpMs <= 0 采取断开策略")
				break
			} else {
				err := conn.SetDeadline(time.Now().Add(time.Millisecond * srv.FrameBrokenDumpMs))
				if err != nil {
					slog.Error("SetDeadline错误", "Error", err.Error())
					break
				}
				_, err = io.ReadAll(conn)
				if err != nil {
					slog.Info(brokenAction, "strategy", fmt.Sprintf("倾倒%dms内的数据", srv.FrameBrokenDumpMs))
					err := conn.SetDeadline(time.Time{})
					if err != nil {
						return
					}
					continue
				}
			}

		}

		if srv.Monitor != nil {
			srv.Monitor.BytesRead <- cr.FrameLength
		}

		if srv.Handler != nil {
			srv.Handler.OnMessage(visitor, cr, 0)
		}

	}
}
