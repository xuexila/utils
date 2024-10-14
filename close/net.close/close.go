package net_close

import "net"

func CloseConn(conn net.Conn) {
	if conn != nil {
		_ = conn.Close()
	}
}
func CloseUdpConn(conn *net.UDPConn) {
	if conn != nil {
		_ = conn.Close()
	}
}
