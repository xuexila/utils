package websocketClose

import (
	"golang.org/x/net/websocket"
)

func Closews(conn *websocket.Conn) {
	if conn != nil {
		_ = conn.Close()
	}
}
