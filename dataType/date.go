package dataType

import (
	"database/sql"
	"database/sql/driver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
	"time"
)

type CustomTime time.Time

func (this *CustomTime) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*this = CustomTime(nullTime.Time)
	return
}

func (this CustomTime) Value() (driver.Value, error) {
	return this, nil
}

func (this CustomTime) GormDataType() string {
	return "time"
}

// GormDBDataType gorm db data type
func (CustomTime) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "timestamp"
	case "mysql":
		return "timestamp"
	case "postgres":
		return "timestamp with time zone"
	case "sqlserver":
		return "timestamp"
	}
	return ""
}

func (this CustomTime) GobEncode() ([]byte, error) {
	return this.GobEncode()
}

func (this *CustomTime) GobDecode(b []byte) error {
	return this.GobDecode(b)
}

func (this CustomTime) MarshalJSON() ([]byte, error) {
	b := []byte{'"'}
	b = append(b, []byte(time.Time(this).Format(time.DateTime))...)
	b = append(b, '"')
	return b, nil
}
func (this *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		*this = CustomTime{}
		return nil
	}
	_t, err := time.ParseInLocation(time.DateTime, s, time.FixedZone("CST", 8*3600))
	*this = CustomTime(_t)
	return err
}
