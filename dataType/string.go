package dataType

import (
	"database/sql"
	"database/sql/driver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type String string

func (this *String) Scan(val interface{}) (err error) {
	nullStr := sql.NullString{}
	err = nullStr.Scan(val)
	*this = String(nullStr.String)
	return
}

func (this String) Value() (driver.Value, error) {
	return string(this), nil
}

func (this String) GormDataType() string {
	return "string"
}

// GormDBDataType gorm db data type
func (String) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "text"
	case "mysql":
		return "LONGTEXT"
	case "postgres":
		return "text"
	case "sqlserver":
		return "VARCHAR(MAX)"
	}
	return ""
}
