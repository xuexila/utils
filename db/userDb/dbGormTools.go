package userDb

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func GormDbDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "longtext"
	case "postgres":
		return "JSONB"
	case "sqlserver":
		return "NVARCHAR(MAX)"
	}
	return ""
}

func GormDbDataValue(d any) (driver.Value, error) {
	if d == nil {
		return nil, nil
	}
	b, err := json.Marshal(d)
	return string(b), err
}

func GormDbScan(val any, dst any) error {
	if val == nil {
		return nil
	}
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", val))
	}
	if len(ba) < 1 {
		return nil
	}
	err := json.Unmarshal(ba, dst)
	return err
}
