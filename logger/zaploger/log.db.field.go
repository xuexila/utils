package zaploger

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

func (this Config) Value() (driver.Value, error) {
	byt, err := json.Marshal(this)
	return string(byt), err
}

func (this *Config) Scan(val any) error {
	if val == nil {
		*this = Config{}
		return nil
	}
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal Config value:", val))
	}
	rd := bytes.NewReader(ba)
	decoder := json.NewDecoder(rd)
	decoder.UseNumber()
	return decoder.Decode(this)
}
