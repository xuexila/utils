package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/url"
)

func main() {
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword("itestor", "wwwitestcom"),
		Host:   "192.168.2.187:5432",
		Path:   "data_house",
	}
	query := dsn.Query()
	query.Set("timezone", "Asia/Shanghai")
	query.Set("search_path", "public")
	dsn.RawQuery = query.Encode()
	dstStr := dsn.String()
	dialector := postgres.New(postgres.Config{
		DSN:                  dstStr,
		PreferSimpleProtocol: true,
	})
	cfg := gorm.Config{
		SkipDefaultTransaction:                   true,
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	db, _ := gorm.Open(dialector, &cfg)

	// 查询当前会话的时区设置
	var timeZone string
	err := db.Raw("SHOW TIMEZONE").Scan(&timeZone).Error
	if err != nil {
		log.Fatal("失败", err)
	}
	fmt.Println("Current session timezone:", timeZone)
}
