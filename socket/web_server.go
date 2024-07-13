package socket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/grammars/easy-go/sugar"
	"log/slog"
	"net/http"
)

type WebServer struct {
	Port            int
	PrintDetail     bool
	Monitor         *Monitor
	upgrader        *websocket.Upgrader
	ReadBufferSize  int
	WriteBufferSize int
}

func (srv *WebServer) Start(ginEngine *gin.Engine) (*gin.Engine, error) {
	slog.Info("WebServer 开始启动")
	if ginEngine == nil {
		ginEngine = gin.Default()
	}
	ginEngine.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{"name": "🐑", "age": 18})
	})
	srv.upgrader = &websocket.Upgrader{
		ReadBufferSize:  sugar.EnsurePositive(srv.ReadBufferSize, 64),
		WriteBufferSize: sugar.EnsurePositive(srv.WriteBufferSize, 64),
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ginEngine.GET("/ws", srv.wsHandler)
	addr := fmt.Sprintf("0.0.0.0:%d", srv.Port)
	err := ginEngine.Run(addr)
	if err != nil {
		return nil, err
	}
	return ginEngine, nil
}

func (srv *WebServer) wsHandler(c *gin.Context) {
	conn, err := srv.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("Error upgrading to websocket:", err)
		return
	}
	slog.Info("Accept new websocket connection")
	defer CloseWebConn(conn)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			slog.Error("Error read message from websocket:", err)
			break
		}
		slog.Info("收到WebSocket发来的消息", "message", message)
		err = conn.WriteMessage(messageType, []byte("俺收到了消息"))
		if err != nil {
			slog.Error("Error write message from websocket:", err)
			break
		}
	}
}
