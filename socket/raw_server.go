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
	Port        int
	PrintDetail bool

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
		if srv.PrintDetail {
			slog.Info("有一个客户端连接我成功了", "RemoteAddr", conn.RemoteAddr())
		}
		go ReadWriteAsServer(conn, srv)
	}
}

func (srv *RawServer[VD]) GetStartTime() time.Time {
	return srv.StartTime
}

type RawVisitorConnection struct {
	conn       net.Conn
	WriteMutex sync.Mutex
}

func (rvc *RawVisitorConnection) RemoteAddr() net.Addr {
	return rvc.conn.RemoteAddr()
}

func (rvc *RawVisitorConnection) Write(b []byte) (n int, err error) {
	return rvc.conn.Write(b)
}

func (rvc *RawVisitorConnection) WriteSafe(b []byte) (n int, err error) {
	rvc.WriteMutex.Lock()
	defer rvc.WriteMutex.Unlock()
	return rvc.Write(b)
}

func (srv *RawServer[VD]) appendVisitor(conn net.Conn) *Visitor[VD] {
	rvc := &RawVisitorConnection{conn: conn}
	visitor := srv.VisitorMap.Append(rvc)
	slog.Info("Accept客户端", "Uid", visitor.Uid, "index", visitor.index, "addr", conn.RemoteAddr())
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
		slog.Info("本帧长度", "FrameLength", cr.FrameLength)

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
