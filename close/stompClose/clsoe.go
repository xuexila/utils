package stompClose

import "github.com/go-stomp/stomp/v3"

func Close(conn *stomp.Conn) {
	if conn != nil {
		_ = conn.Disconnect()
	}
}
