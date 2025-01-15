package userDb

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/helays/utils/dataType"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// GormDbDataTypeJSON 弃用函数，推荐使用 dataType.JsonDbDataType(db, field)
// deprecated
func GormDbDataTypeJSON(db *gorm.DB, field *schema.Field) string {
	return dataType.JsonDbDataType(db, field)
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
