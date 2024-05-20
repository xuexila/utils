package redigo

import "github.com/gomodule/redigo/redis"

// CloseRedisConn github.com/gomodule/redigo/redis
func CloseRedisConn(conn redis.Conn) {
	if conn != nil {
		_ = conn.Close()
	}
}
