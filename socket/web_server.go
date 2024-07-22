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
	"time"
)

type WebServer[VD any] struct {
	Port            int
	PrintDetail     bool
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

func (srv *WebServer[VD]) appendVisitor(conn *websocket.Conn) *Visitor[VD] {
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

func (srv *WebServer[VD]) removeVisitor(visitorUid uint64) {
	srv.VisitorMap.Remove(visitorUid)
	slog.Info("Remove客户端", "visitorUid", visitorUid)
	if srv.PrintDetail {
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
	defer srv.removeVisitor(visitor.uid)
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
		if srv.Monitor != nil {
			srv.Monitor.BytesRead <- len(message)
		}
		if srv.PrintDetail {
			messageLen := len(message)
			if messageLen < 24 {
				slog.Info("收到WebSocket发来的消息", "message", fmt.Sprintf("%x", message))
			} else {
				partMessage := message[:24]
				slog.Info("收到WebSocket发来的消息", "message", fmt.Sprintf("%x len=%d", partMessage, messageLen))
			}
		}
		if srv.Handler != nil {
			srv.Handler.OnMessage(visitor, &message)
		}
		resp := []byte("俺收到了消息")
		err = conn.WriteMessage(messageType, resp)
		if err != nil {
			slog.Error("Error write message from websocket:", err)
			break
		}
		srv.Monitor.BytesWrite <- len(resp)
	}
}
