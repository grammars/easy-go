package socket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/url"
	"time"
)

type WebClient struct {
	Host        string
	Port        int
	WsPath      string
	Name        string
	PrintDetail bool

	Monitor *Monitor
}

func (cli *WebClient) Start() {
	addr := fmt.Sprintf("%s:%d", cli.Host, cli.Port)
	slog.Info("WebClient 开始启动", "Name", cli.Name, "addr", addr)
	u := url.URL{Scheme: "ws", Host: addr, Path: GetWsPath(cli.WsPath)}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		slog.Error("连接失败", "url", u.String(), "err", err.Error())
		return
	}
	if cli.Monitor != nil {
		cli.Monitor.ValidNum <- 1
	}
	defer func() {
		if cli.Monitor != nil {
			cli.Monitor.InvalidNum <- 1
		}
	}()
	defer CloseWebConn(conn)

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				slog.Error("读取消息失败", "Name", cli.Name, "err", err.Error())
				return
			}
			if cli.Monitor != nil {
				cli.Monitor.BytesRead <- len(message)
			}
			if cli.PrintDetail {
				slog.Info("读取到了:", "message", string(message))
			}
		}
	}()

	var i = 0
	for {
		i++
		message := fmt.Sprintf("Hi Web-Server-(%d)", i)
		resp := []byte(message)
		err := conn.WriteMessage(websocket.TextMessage, resp)
		if err != nil {
			slog.Error("写消息失败", "Name", cli.Name, "err", err.Error())
			return
		}
		if cli.Monitor != nil {
			cli.Monitor.BytesWrite <- len(resp)
		}
		time.Sleep(200 * time.Millisecond)
		if i >= 5*3600 {
			break
		}
	}
}
