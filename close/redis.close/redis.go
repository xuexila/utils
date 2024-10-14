package redis_close

import rds "github.com/redis/go-redis/v9"

func CloseRdsConn(conn *rds.Conn) {
	if conn != nil {
		_ = conn.Close()
	}
}

func CloseRedisUniversalClient(rdb rds.UniversalClient) {
	if rdb != nil {
		_ = rdb.Close()
	}
}
