package config

// interval 定义
const (
	IntervalSecond = "SECOND"
	IntervalMinute = "MINUTE"
	IntervalHour   = "HOUR"
	IntervalDay    = "DAY"
	IntervalWeek   = "WEEK"
	IntervalMonth  = "MONTH"
	IntervalYear   = "YEAR"

	IntervalSecondLabel = "秒"
	IntervalMinuteLabel = "分"
	IntervalHourLabel   = "时"
	IntervalDayLabel    = "天"
	IntervalWeekLabel   = "周"
	IntervalMonthLabel  = "月"
	IntervalYearLabel   = "年"
)

type Interval struct {
	Key   int    `json:"key" yaml:"key"`     // key
	Val   string `json:"val" yaml:"val"`     // val
	Label string `json:"Label" yaml:"Label"` //单位
}

var IntervalLists = []Interval{
	{
		Key:   1,
		Val:   IntervalSecond,
		Label: IntervalSecondLabel,
	},
	{
		Key:   60,
		Val:   IntervalMinute,
		Label: IntervalMinuteLabel,
	},
	{
		Key:   3600,
		Val:   IntervalHour,
		Label: IntervalHourLabel,
	},
	{
		Key:   86400,
		Val:   IntervalDay,
		Label: IntervalDayLabel,
	},
	{
		Key:   604800,
		Val:   IntervalWeek,
		Label: IntervalWeekLabel,
	},
	{
		Key:   2592000,
		Val:   IntervalMonth,
		Label: IntervalMonthLabel,
	},
	{
		Key:   31536000,
		Val:   IntervalYear,
		Label: IntervalYearLabel,
	},
}

var IntervalKeyMap = make(map[int]Interval)
var IntervalKeyValMap = make(map[string]Interval)
var IntervalKeyLabelMap = make(map[string]Interval)

// 自动初始化
func init() {
	for _, v := range IntervalLists {
		IntervalKeyMap[v.Key] = v
		IntervalKeyValMap[v.Val] = v
		IntervalKeyLabelMap[v.Label] = v
	}
}

// 关系数据库
const (
	DbTypeMysql      = "mysql"
	DbTypePostgres   = "postgres"
	DbTypePostgresql = "postgresql"
	DbTypePg         = "pg"
	DbTypeSqlite     = "sqlite"
	DbTypeMssql      = "mssql"
	DbTypeOracle     = "oracle"
	DbTypeSqlserver  = "sqlserver"
)

// kv数据库
const (
	DbTypeMongo = "mongo"
	DbTypeRedis = "redis"
)

// 搜索数据库
const (
	DbTypeEs = "es" // es存储
)

// 消息队列
const (
	QueueTypeKafka    = "kafka" // kafka消息队列
	QueueTypeRabbit   = "rabbit"
	QueueTypeRocketmq = "rocketmq"
	QueueTypeRabbitmq = "rabbitmq"
)

// 文件存储
const (
	FileTypeFtp   = "ftp"
	FileTypeSftp  = "sftp"
	FileTypeLocal = "local"
	FileTypeOss   = "oss"
	FileTypeMinio = "minio"
	FileType      = "hdfs"
)

// 集群类型
const (
	ClusterEtcd      = "etcd"
	ClusterNacos     = "nacos"
	ClusterZookeeper = "zookeeper"
)
