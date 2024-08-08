package socket

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/grammars/easy-go/best"
	"github.com/grammars/easy-go/sugar"
	"io"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"
)

type WebServer[VD any] struct {
	Port            int
	WsPath          string
	ReadBufferSize  int
	WriteBufferSize int
	TLS             bool
	CrtFile         string
	KeyFile         string

	upgrader *websocket.Upgrader

	StartTime  time.Time
	Monitor    *Monitor
	Handler    VisitorServerHandler[VD]
	VisitorMap *VisitorMap[VD]
}

func (srv *WebServer[VD]) StartDefault() {
	_, err := srv.Start(nil)
	if err != nil {
		slog.Error("启动WebServer失败", "Error", err)
		return
	}
}

func (srv *WebServer[VD]) Start(ginEngine *gin.Engine) (*gin.Engine, error) {
	slog.Info("WebServer 开始启动")
	srv.StartTime = time.Now()
	if ginEngine == nil {
		ginEngine = gin.Default()
	}
	ginEngine.GET("/status", func(c *gin.Context) {
		if srv.Monitor != nil {
			c.JSON(200, best.SuccessResult("已开启统计", srv.Monitor.Stat.ToMap()))
		} else {
			c.JSON(200, best.FailResult("未开启统计"))
		}
	})
	srv.VisitorMap = CreateVisitorMap[VD](srv)
	srv.upgrader = &websocket.Upgrader{
		ReadBufferSize:  sugar.EnsurePositive(srv.ReadBufferSize, 64),
		WriteBufferSize: sugar.EnsurePositive(srv.WriteBufferSize, 64),
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ginEngine.GET(GetWsPath(srv.WsPath), srv.wsHandler)
	if srv.Port <= 0 {
		srv.Port = sugar.ReturnIf(srv.TLS, 443, 80)
	}
	addr := fmt.Sprintf("0.0.0.0:%d", srv.Port)
	slog.Info("WebServer 监听", "addr", addr)
	var err error
	if srv.TLS {
		err = ginEngine.RunTLS(addr, srv.CrtFile, srv.KeyFile)
	} else {
		err = ginEngine.Run(addr)
	}
	if err != nil {
		return nil, err
	}
	return ginEngine, nil
}

func (srv *WebServer[VD]) GetStartTime() time.Time {
	return srv.StartTime
}

type WebVisitorConnection[VD any] struct {
	conn       *websocket.Conn
	WriteMutex sync.Mutex
	srv        *WebServer[VD]
}

func (wvc *WebVisitorConnection[VD]) RemoteAddr() net.Addr {
	return wvc.conn.RemoteAddr()
}

func (wvc *WebVisitorConnection[VD]) Write(b []byte) (int, error) {
	e := wvc.conn.WriteMessage(websocket.BinaryMessage, b)
	nBytes := len(b)
	if wvc.srv.Monitor != nil {
		wvc.srv.Monitor.BytesWrite <- nBytes
	}
	return nBytes, e
}

func (wvc *WebVisitorConnection[VD]) WriteSafe(b []byte) (int, error) {
	wvc.WriteMutex.Lock()
	defer wvc.WriteMutex.Unlock()
	return wvc.Write(b)
}

func (srv *WebServer[VD]) appendVisitor(conn *websocket.Conn) *Visitor[VD] {
	wvc := &WebVisitorConnection[VD]{conn: conn, srv: srv}
	visitor := srv.VisitorMap.Append(wvc)
	if LogLevel <= 0 {
		slog.Info("web_server.go Accept访问者", "uid", visitor.Uid, "index", visitor.index, "addr", conn.RemoteAddr())
		srv.VisitorMap.Print()
	}
	if srv.Monitor != nil {
		srv.Monitor.ValidNum <- 1
	}
	return visitor
}

func (srv *WebServer[VD]) removeVisitor(visitorUid uint64) {
	srv.VisitorMap.Remove(visitorUid)
	if LogLevel <= 0 {
		slog.Info("web_server.go Remove访问者", "visitorUid", visitorUid)
		srv.VisitorMap.Print()
	}
	srv.Monitor.InvalidNum <- 1
}

func (srv *WebServer[VD]) wsHandler(c *gin.Context) {
	conn, err := srv.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("Error upgrading to websocket:", err)
		return
	}
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
	defer CloseWebConn(conn)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil && err != io.EOF {
			var wce *websocket.CloseError
			if errors.As(err, &wce) {
				slog.Warn("客户端已断开(websocket.CloseError)", "Code", wce.Code, "Text", wce.Text)
			} else {
				var noe *net.OpError
				if errors.As(err, &noe) {
					slog.Warn("客户端已断开(net.OpError)", "Op", noe.Op, "Error", noe.Error())
				} else {
					slog.Error("Error read message from websocket:", err)
				}
			}
			break
		}
		messageLen := len(message)
		if srv.Monitor != nil {
			srv.Monitor.BytesRead <- messageLen
		}
		if LogLevel <= 0 {
			slog.Info("web_server.go 收到WebSocket发来的消息", "message", SprintWebSocketFrame(messageType, message, 24))
		}
		if srv.Handler != nil {
			msg := CodecResult{FrameLength: messageLen, BodyBytes: message}
			srv.Handler.OnMessage(visitor, msg, messageType)
		}
	}
}
