package dataType

import (
	"database/sql"
	"database/sql/driver"
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

func (this *CustomTime) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	this.Time = nullTime.Time
	return
}

func (this CustomTime) Value() (driver.Value, error) {
	return this.Time, nil
}

func (this CustomTime) GormDataType() string {
	return "time"
}

func (this CustomTime) GobEncode() ([]byte, error) {
	return this.Time.GobEncode()
}

func (this *CustomTime) GobDecode(b []byte) error {
	return (this.Time).GobDecode(b)
}

func (this CustomTime) MarshalJSON() ([]byte, error) {
	b := []byte{'"'}
	b = append(b, []byte(this.Time.Format(time.DateTime))...)
	b = append(b, '"')
	return b, nil
}
func (this *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		this.Time = time.Time{}
		return nil
	}
	this.Time, err = time.ParseInLocation(time.DateTime, s, time.FixedZone("CST", 8*3600))
	return err
}
