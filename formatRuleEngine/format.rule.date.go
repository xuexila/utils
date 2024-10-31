package formatRuleEngine

import (
	"github.com/araddon/dateparse"
	"github.com/helays/utils/tools"
	"time"
)

// 时间格式化
func (this FormatRule) dateFormat(_src any) (string, error) {
	var (
		t   time.Time
		err error
		src = tools.Any2string(_src)
	)
	// 首先尝试使用 https://github.com/araddon/dateparse 库
	t, err = dateparse.ParseLocal(src)
	if err == nil {
		if this.OutputRule != "" {
			return t.Format(this.OutputRule), nil
		}
		return t.Format(time.DateTime), err
	}
	for _, format := range this.InputRules {
		if format == "timestamp" {
			t, err = tools.AutoDetectTimestampString(src)
		} else {
			t, err = time.Parse(format, src)
		}
		if err != nil {
			continue
		}
		if this.OutputRule != "" {
			return t.Format(this.OutputRule), nil
		}
		return t.Format(time.DateTime), err
	}
	return src, nil
}
