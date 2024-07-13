package socket

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/grammars/easy-go/sugar"
	"io"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

type WebServer struct {
	Port            int
	PrintDetail     bool
	WsPath          string
	ReadBufferSize  int
	WriteBufferSize int

	StartTime      time.Time
	Monitor        *Monitor
	upgrader       *websocket.Upgrader
	visitorHistory uint64
	visitorMap     map[uint64]*WebVisitor
}

type WebVisitor struct {
	index uint64
	uid   uint64
	conn  *websocket.Conn
}

func (srv *WebServer) StartDefault() {
	_, err := srv.Start(nil)
	if err != nil {
		slog.Error("启动WebServer失败", "Error", err)
		return
	}
}

func (srv *WebServer) Start(ginEngine *gin.Engine) (*gin.Engine, error) {
	slog.Info("WebServer 开始启动")
	srv.StartTime = time.Now()
	if ginEngine == nil {
		ginEngine = gin.Default()
	}
	ginEngine.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{"name": "🐑", "age": 18})
	})
	srv.visitorMap = make(map[uint64]*WebVisitor)
	srv.upgrader = &websocket.Upgrader{
		ReadBufferSize:  sugar.EnsurePositive(srv.ReadBufferSize, 64),
		WriteBufferSize: sugar.EnsurePositive(srv.WriteBufferSize, 64),
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ginEngine.GET(GetWsPath(srv.WsPath), srv.wsHandler)
	addr := fmt.Sprintf("0.0.0.0:%d", srv.Port)
	err := ginEngine.Run(addr)
	if err != nil {
		return nil, err
	}
	return ginEngine, nil
}

func (srv *WebServer) printVisitorMap() {
	for k, v := range srv.visitorMap {
		slog.Info("打印visitorMap", "uid", k, "index", v.index, "addr", v.conn.RemoteAddr())
	}
}

func (srv *WebServer) appendVisitor(conn *websocket.Conn) *WebVisitor {
	visitor := WebVisitor{conn: conn}
	visitor.index = atomic.AddUint64(&srv.visitorHistory, 1)
	visitor.uid = uint64(srv.StartTime.UnixMilli()) + visitor.index
	srv.visitorMap[visitor.uid] = &visitor
	slog.Info("Accept客户端", "uid", visitor.uid, "index", visitor.index, "addr", conn.RemoteAddr())
	if srv.PrintDetail {
		srv.printVisitorMap()
	}
	if srv.Monitor != nil {
		srv.Monitor.ValidNum <- 1
	}
	return &visitor
}

func (srv *WebServer) removeVisitor(visitorUid uint64) {
	delete(srv.visitorMap, visitorUid)
	slog.Info("Remove客户端", "visitorUid", visitorUid)
	if srv.PrintDetail {
		srv.printVisitorMap()
	}
	srv.Monitor.InvalidNum <- 1
}

func (srv *WebServer) wsHandler(c *gin.Context) {
	conn, err := srv.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("Error upgrading to websocket:", err)
		return
	}
	visitor := srv.appendVisitor(conn)
	defer srv.removeVisitor(visitor.uid)
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
			slog.Info("收到WebSocket发来的消息", "message", message)
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
