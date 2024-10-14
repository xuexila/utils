package websocket_close

import (
	"golang.org/x/net/websocket"
)

func Closews(conn *websocket.Conn) {
	if conn != nil {
		_ = conn.Close()
	}
}
