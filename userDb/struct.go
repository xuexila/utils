package userDb

type Dbbase struct {
	DbType       string   `ini:"db_type" yaml:"db_type" json:"db_type"` // 数据库类型 mysql/pg
	Host         []string `ini:"host" yaml:"host" json:"host"`
	User         string   `ini:"user" yaml:"user" json:"user"`
	Pwd          string   `ini:"pwd" yaml:"pwd" json:"pwd"`
	Dbname       string   `ini:"dbname" yaml:"dbname" json:"dbname"`
	MaxIdleConns int      `ini:"max_idle_conns" yaml:"maxIdleConns" json:"maxIdleConns"`
	MaxOpenConns int      `ini:"max_open_conns" yaml:"maxOpenConns" json:"maxOpenConns"`
}
