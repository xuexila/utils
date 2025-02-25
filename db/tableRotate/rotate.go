package tableRotate

import (
	"context"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/helays/utils/config"
	"github.com/helays/utils/logger/ulogs"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"strings"
	"time"
)

// TableRotate db自动轮转配置
type TableRotate struct {
	Enable                  bool          `json:"enable" yaml:"enable" ini:"enable"` // 是否启用自动轮转
	Duration                time.Duration `json:"duration" yaml:"duration" ini:"duration"`
	Crontab                 string        `json:"crontab" yaml:"crontab" ini:"crontab"`                                                          // crontab 表达式 ,定时器和crontab二选一
	SplitTable              bool          `json:"split_table" yaml:"split_table" ini:"split_table"`                                              // 是否开启按天切分日志 ，开启后，自动回收数据 只会看表的保留数量，不开启，就看数据保留时长
	MaxTableRetention       int           `json:"max_table_retention" yaml:"max_table_retention" ini:"max_table_retention"`                      // 最大保留天数 -1 不限制
	DataRetentionPeriod     int           `json:"data_retention_period" yaml:"data_retention_period" ini:"data_retention_period"`                // 数据保留时长 -1 不限制
	DataRetentionPeriodUnit string        `json:"data_retention_period_unit" yaml:"data_retention_period_unit" ini:"data_retention_period_unit"` // 数据保留时间单位 支持 second minute hour day month year
	FilterField             string        `json:"filter_field" yaml:"filter_field" ini:"filter_field"`                                           // 过滤字段 默认create_time
	tx                      *gorm.DB
	tableName               string
}

// AddTask 添加自动轮转任务
func (this TableRotate) AddTask(ctx context.Context, tx *gorm.DB, tableName string) {
	if !this.Enable {
		return
	}
	if this.Duration < 0 && this.Crontab == "" {
		return
	}
	ulogs.Log("数据库", tx.Dialector.Name(), tableName, "配置自动轮转")
	this.tx = tx
	this.tableName = tableName
	if this.Crontab != "" {
		c := cron.New()
		if _, err := c.AddFunc(this.Crontab, this.run); err != nil {
			ulogs.Error("添加自动轮转任务失败", "表", tableName, "定时", this.Crontab)
			return
		}
		c.Start()
		return
	}
	go func() {
		ulogs.Log("数据库", tx.Dialector.Name(), tableName, "配置自动轮转", "间隔", this.Duration)
		tck := time.NewTicker(this.Duration)
		defer tck.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tck.C:
				this.run()
			}
		}
	}()
}

func (this *TableRotate) run() {
	if this.SplitTable {
		this.runSplitTable()
		return
	}
	this.runRotateTableData()
}

// 分表
func (this *TableRotate) runSplitTable() {
	if this.MaxTableRetention <= 0 {
		return
	}
	// 然后进行分表切割
}

// 回收表数据
func (this *TableRotate) runRotateTableData() {
	if this.DataRetentionPeriod <= 0 {
		return
	}
	var queryVal clause.Expr
	unit := strings.ToUpper(this.DataRetentionPeriodUnit)
	retentionPeriod := strconv.Itoa(this.DataRetentionPeriod)
	switch this.tx.Dialector.Name() {
	case config.DbTypeMysql:
		queryVal = clause.Expr{
			SQL:                "NOW() - INTERVAL '? ?'",
			Vars:               []any{clause.Column{Name: retentionPeriod, Raw: true}, clause.Column{Name: unit, Raw: true}},
			WithoutParentheses: false,
		}
	case config.DbTypePostgres:
		queryVal = clause.Expr{
			SQL:                "NOW() - INTERVAL '? ?'",
			Vars:               []any{clause.Column{Name: retentionPeriod, Raw: true}, clause.Column{Name: unit, Raw: true}},
			WithoutParentheses: false,
		}
	}
	err := this.tx.Table(this.tableName).Where(this.FilterField+" < ?", queryVal).Delete(nil).Error
	if err != nil {
		switch _err := err.(type) {
		case *pgconn.PgError:
			if _err.Code == "42P01" {
				return
			}
		case *mysql.MySQLError:

		default:
			fmt.Println("fadsf", _err)
		}
		ulogs.Error("自动轮转表，回收表数据失败", this.tableName, "过滤字段", this.FilterField, "条件", retentionPeriod, unit, err)
	} else {
		ulogs.Log("自动轮转表", this.tableName, "回收表数据成功", "过滤字段", this.FilterField, "条件", retentionPeriod, unit)
	}

}
