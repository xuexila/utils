package utils

import (
	"database/sql"
	"database/sql/driver"
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

func (date *CustomTime) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	date.Time = nullTime.Time
	return
}

func (date CustomTime) Value() (driver.Value, error) {
	return date.Time, nil
}

// GormDataType gorm common data type
func (date CustomTime) GormDataType() string {
	return "date"
}

func (date CustomTime) GobEncode() ([]byte, error) {
	return time.Time(date.Time).GobEncode()
}

func (date *CustomTime) GobDecode(b []byte) error {
	return (date.Time).GobDecode(b)
}

func (date CustomTime) MarshalJSON() ([]byte, error) {
	return date.Time.MarshalJSON()
}
func (this *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		this.Time = time.Time{}
		return nil
	}
	this.Time, err = time.ParseInLocation("2006-01-02 15:04:05", s, time.FixedZone("CST", 8*3600))
	return err
}
