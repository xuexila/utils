package dataType

import (
	"context"
	"database/sql/driver"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type Array []any

// Value return json value, implement driver.Valuer interface
func (m Array) Value() (driver.Value, error) {
	return arrayValue(m)
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *Array) Scan(val interface{}) error {
	return arrayScan(m, val)
}

// GormDataType gorm common data type
func (m Array) GormDataType() string {
	return "jsonmap"
}

// GormDBDataType gorm db data type
func (Array) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return jsonDbDataType(db, field)
}
func (jm Array) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return arrayGormValue(jm, db)
}

type StringArray []string

func (m StringArray) Value() (driver.Value, error) {
	b, err := arrayValue(m)
	return b, err
}
func (m *StringArray) Scan(val interface{}) error {
	return arrayScan(m, val)
}

func (StringArray) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return jsonDbDataType(db, field)
}

func (jm StringArray) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return arrayGormValue(jm, db)
}
