package dataType

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"strings"
)

// JSONMap defined JSON data type, need to implements driver.Valuer, sql.Scanner interface
type JSONMap map[string]any

// Value return json value, implement driver.Valuer interface
func (m JSONMap) Value() (driver.Value, error) {
	return DriverValueWithJson(m)
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *JSONMap) Scan(val any) error {
	return DriverScanWithJson(val, m) // 这里暂时先用这个版本
	//if val == nil {
	//	*m = make(JSONMap)
	//	return nil
	//}
	//var ba []byte
	//switch v := val.(type) {
	//case []byte:
	//	ba = v
	//case string:
	//	ba = []byte(v)
	//default:
	//	return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", val))
	//}
	//t := map[string]any{}
	//rd := bytes.NewReader(ba)
	//decoder := json.NewDecoder(rd)
	//decoder.UseNumber()
	//err := decoder.Decode(&t)
	//*m = t
	//return err
}

// MarshalJSON to output non base64 encoded []byte
func (m JSONMap) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	t := (map[string]any)(m)
	return json.Marshal(t)
}

// UnmarshalJSON to deserialize []byte
func (m *JSONMap) UnmarshalJSON(b []byte) error {
	rd := bytes.NewReader(b)
	decoder := json.NewDecoder(rd)
	decoder.UseNumber()
	return decoder.Decode(m)
}

// GormDataType gorm common data type
func (m JSONMap) GormDataType() string {
	return "jsonmap"
}

// GormDBDataType gorm db data type
func (JSONMap) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return JsonDbDataType(db, field)
}

func (jm JSONMap) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := jm.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") && CheckVersionSupportsJSON(v.ServerVersion) {
			fmt.Println(v.ServerVersion)
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
