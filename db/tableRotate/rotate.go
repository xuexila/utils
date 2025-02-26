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
	"sort"
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
	ulogs.Log("【表自动轮转配置】", "数据库", tx.Dialector.Name(), tx.Migrator().CurrentDatabase(), tableName)
	ulogs.Log("【表自动轮转配置】", "周期策略", this.Crontab, this.Duration)
	if this.SplitTable {
		ulogs.Log("【表自动轮转配置】", "回收策略：", "分表", "最大保留数量", this.MaxTableRetention)
	} else {
		ulogs.Log("【表自动轮转配置】", "回收策略：", "数据", "数据保留时长", this.DataRetentionPeriod, this.DataRetentionPeriodUnit)
	}
	this.tx = tx
	this.tableName = tableName
	if this.Crontab != "" {
		go this.toCrontab(ctx)
		return
	}
	go this.toTicker(ctx)
}

// 通过 crontab方式运行
func (this *TableRotate) toCrontab(ctx context.Context) {
	c := cron.New()
	eid, err := c.AddFunc(this.Crontab, this.run)
	if err != nil {
		ulogs.Error("添加自动轮转任务失败", "表", this.tableName, "定时", this.Crontab)
		return
	}
	c.Start()
	go func() {
		<-ctx.Done()      // 等待上下文取消
		c.Remove(eid)     // 移除任务
		<-c.Stop().Done() // 停止 cron 调度器
		ulogs.Log("【表自动轮转配置终止】", "crontab", "数据库", this.tx.Dialector.Name(), this.tx.Migrator().CurrentDatabase(), this.tableName)
	}()
}

// 通过 定时器方式运行
func (this *TableRotate) toTicker(ctx context.Context) {
	tck := time.NewTicker(this.Duration)
	defer tck.Stop()
	for {
		select {
		case <-ctx.Done():
			ulogs.Log("【表自动轮转配置终止】", "定时器", "数据库", this.tx.Dialector.Name(), this.tx.Migrator().CurrentDatabase(), this.tableName)
			return
		case <-tck.C:
			this.run()
		}
	}
}

func (this *TableRotate) run() {
	if this.SplitTable {
		this.runSplitTable()
		return
	}
	this.runRotateTableData()
}

const dateFormat = "20060102150405"

// 分表
func (this *TableRotate) runSplitTable() {
	if this.MaxTableRetention <= 0 {
		return
	}

	// 然后进行分表切割

	newTableName := this.tableName + "_" + time.Now().Format(dateFormat)
	err := this.tx.Transaction(func(tx *gorm.DB) error {
		err := tx.Migrator().RenameTable(this.tableName, newTableName)
		if err != nil {
			return fmt.Errorf("修改表名失败 %s to %s :%s", this.tableName, newTableName, err.Error())
		}
		switch tx.Dialector.Name() {
		case config.DbTypePostgres:
			// 创建新表
			err = tx.Debug().Exec("CREATE TABLE ? (LIKE ? INCLUDING ALL)", clause.Table{Name: this.tableName}, clause.Table{Name: newTableName}).Error
		case config.DbTypeMysql:
			err = tx.Debug().Exec("CREATE TABLE ? LIKE ?", clause.Table{Name: this.tableName}, clause.Table{Name: newTableName}).Error
		}
		if err != nil {
			return fmt.Errorf("创建表失败 %s :%s", this.tableName, err.Error())
		}
		return nil
	})
	if err != nil {
		ulogs.Error("自动轮转表，修改表名失败", this.tableName, "新表名", newTableName, err)
	} else {
		ulogs.Log("自动轮转表", this.tableName, "修改表名成功", "新表名", newTableName)
	}
	// 查询以 this.tableName开头的表名
	var tables []string
	switch this.tx.Dialector.Name() {
	case config.DbTypePostgres:
		// 还要查询当前的搜索模式
		// 获取当前搜索模式
		var searchPath string
		if err = this.tx.Raw("SHOW search_path").Scan(&searchPath).Error; err != nil {
			ulogs.Error("自动轮转表，获取当前搜索模式失败", err)
			return
		}
		// 默认搜索模式是第一个模式
		currentSchema := "public" // 默认值
		if len(searchPath) > 0 {
			currentSchema = strings.Split(searchPath, ",")[0] // 取第一个模式
			currentSchema = strings.TrimSpace(currentSchema)  // 去除空格
		}
		tx := this.tx.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema like ? and table_name LIKE ?", currentSchema, this.tableName+"%")
		err = tx.Scan(&tables).Error
	case config.DbTypeMysql:
		currentDataBase := this.tx.Migrator().CurrentDatabase()
		tx := this.tx.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema like ? and table_name LIKE ?", currentDataBase, this.tableName+"%")
		err = tx.Scan(&tables).Error
	}
	if err != nil {
		ulogs.Error("自动轮转表，查询表名失败", err)
		return
	}

	var ts byCreateTime
	// 获取表名清单，满足tablename_20060102150405 格式的单独计算，然后根据时间倒叙，保留最大数量的表
	for _, tableName := range tables {
		// 根据去除 tableName中的前缀
		createTime := ""
		currentTable := this.tableName + "_"
		if strings.HasPrefix(tableName, currentTable) {
			createTime = tableName[len(currentTable):]
		}
		if t, err := time.Parse(dateFormat, createTime); err == nil {
			ts = append(ts, tableSplit{
				tableName:  tableName,
				createTime: t,
			})
		}
	}
	sort.Sort(ts)
	if len(ts) < this.MaxTableRetention {
		return
	}
	// 删除多余的表
	for _, item := range ts[this.MaxTableRetention:] {
		err = this.tx.Migrator().DropTable(item.tableName)
		if err != nil {
			ulogs.Error("自动轮转表，删除表失败", err)
		} else {
			ulogs.Log("自动轮转表", this.tableName, "删除表成功", "表名", item.tableName)
		}
	}

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
