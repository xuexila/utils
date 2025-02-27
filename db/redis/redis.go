package redis

import (
	"context"
	"fmt"
	"github.com/helays/utils/logger/ulogs"
	"github.com/redis/go-redis/v9"
)

// 配置说明

type Rediscfg struct {
	MasterName       string   `json:"master_name" yaml:"master_name" ini:"master_name"`                    // 指定这个 MasterName ，就是FailoverClient哨兵模式，
	Addrs            []string `json:"addrs" yaml:"addrs" ini:"addrs,omitempty"`                            // 如果这个有两个以及上，就是集群模式
	SentinelAddrs    []string `json:"sentinel_addrs" yaml:"sentinel_addrs" ini:"sentinel_addrs,omitempty"` // 哨兵节点地址列表
	ClientName       string   `json:"client_name" yaml:"client_name" ini:"client_name"`                    // 每个Node节点的每个网络连接配置
	User             string   `json:"user" yaml:"user" ini:"user"`
	Password         string   `json:"password" yaml:"password" ini:"password"`
	SentinelUsername string   `json:"sentinel_username" yaml:"sentinel_username" ini:"sentinel_username"` // 用于ACL认证的用户名
	// Sentinel中 `requirepass<password>` 的密码配置
	// 如果同时提供了 `SentinelUsername` ，则启用ACL认证
	SentinelPassword string `json:"sentinel_password" yaml:"sentinel_password" ini:"sentinel_password"`
	Db               int    `json:"db" yaml:"db" ini:"db"` // 默认数据库
	PoolSize         int    `json:"pool_size" yaml:"pool_size" ini:"pool_size"`
	// 具体来说，当 DisableIndentity 设置为 true 时，它会阻止客户端在建立连接时自动发送命令来设置自己的标识信息。
	// 这通常涉及到通过 CLIENT SETINFO LIBRARY 或类似的命令向 Redis 服务器报告客户端库的名称和版本等信息。
	// 在某些情况下，这可能会导致一些问题，例如，当客户端库不支持这些命令时，或者当应用程序需要控制客户端标识信息的设置方式时。
	DisableIndentity         bool   `json:"disable_identity" yaml:"disable_identity" ini:"disable_identity"`                                     //  是否禁用在连接时设置客户端库标识的行为
	IdentitySuffix           string `json:"identity_suffix" yaml:"identity_suffix" ini:"identity_suffix"`                                        // 默认为空, 用于在客户端标识信息中添加后缀
	EnableSetDbBeforeConnect bool   `json:"enable_set_db_before_connect" yaml:"enable_set_db_before_connect" ini:"enable_set_db_before_connect"` // 在连接前启用设置db编号
	EnableSetDbAfterConnect  bool   `json:"enable_set_db_after_connect" yaml:"enable_set_db_after_connect" ini:"enable_set_db_after_connect"`    // 在连接后启用设置db编号
	EnableAuthOnConnect      bool   `json:"enable_auth_on_connect" yaml:"enable_auth_on_connect" ini:"enable_auth_on_connect"`                   // 在 自定义连接中使用auth ,跟下面的 OnConnect效果一样
	OnConnect                bool   `json:"on_connect" yaml:"on_connect" ini:"on_connect"`                                                       // 主要是在云组件ctg cache 的时候，才需要这个，其他情况一般不需要
	CustomScan               bool   `ini:"custom_scan" yaml:"custom_scan" json:"custom_scan"`                                                    // 系统中使用scan 扫描的时候，云组件可能需要用这个
	EnableCheckOnInit        bool   `ini:"enable_check_on_init" yaml:"enable_check_on_init" json:"enable_check_on_init"`                         // 是否在初始化的时候启用 ping测试
}

// NewUniversalClient 创建一个通用的 Redis 客户端
func (this Rediscfg) NewUniversalClient() (redis.UniversalClient, error) {
	c := redis.UniversalOptions{
		Addrs:            this.Addrs,
		ClientName:       this.ClientName,
		Username:         this.User,
		Password:         this.Password,
		SentinelUsername: this.SentinelUsername,
		SentinelPassword: this.SentinelPassword,
		MasterName:       this.MasterName,
		DisableIndentity: this.DisableIndentity,
		IdentitySuffix:   this.IdentitySuffix,
		PoolSize:         this.PoolSize,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			var err error
			if this.EnableAuthOnConnect || this.OnConnect {
				auth := cn.Auth(ctx, this.Password)
				result, _ := auth.Result()
				ulogs.Log("redis 二次认证", "结果", result, auth.Err())
				if err = auth.Err(); err != nil {
					return err
				}
			}
			if this.EnableSetDbAfterConnect {
				// 在连接的时候，设置库
				return cn.Select(ctx, this.Db).Err()
			}
			return nil
		},
	}

	if this.EnableSetDbBeforeConnect {
		c.DB = this.Db
	}
	ulogs.Log("redis连接参数", this.Addrs, "库编号", this.Db, "二次认证", this.OnConnect || this.EnableAuthOnConnect, "set lib", this.DisableIndentity)
	rdb := redis.NewUniversalClient(&c)
	if this.EnableCheckOnInit {
		status := rdb.Ping(context.Background())
		err := status.Err()
		if err != nil {
			return nil, fmt.Errorf("redis指令【%s】执行失败：%s", status.String(), err.Error())
		}
	}
	ulogs.Log("redis连接成功", this.Addrs, "库编号", this.Db, "二次认证", this.OnConnect || this.EnableAuthOnConnect)
	return rdb, nil
}
