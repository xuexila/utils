package formatRuleEngine

type FormatRule struct {
	FormatType string   `yaml:"format_type" json:"format_type"` // 格式化类型
	InputRules []string `yaml:"input_rules" json:"input_rules"` // 识别格式化规则
	OutputRule string   `yaml:"out_rule" json:"output_rule"`    //  输出格式规则
}

// Format 格式化
func (this FormatRule) Format(src any) (any, error) {
	switch this.FormatType {
	case "date_format":
		return this.dateFormat(src)
	}
	return src, nil
}
