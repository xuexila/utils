package excelTools

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
	"strings"
)

// CloseExcel 关闭excel文件
func CloseExcel(f *excelize.File) {
	if f != nil {
		_ = f.Close()
	}
}

// ReadExcelRow 读取excel 指定或者每个sheet的数据，并通过回调函数方式进行处理
// excelFile *excelize.File excel文件
// sheetName []string sheet名称集合，如果不填就是读取所有sheet
// sheetIndex []int sheet索引集合，如果不填就是读取所有sheet，注意，这里需要从0开始
// dataRow int 数据行数，从第几行开始读取有效数据
// dataCol int 数据列数，保证每行有效数据有多少列
// call func(sheetName string,tmp map[string]any) error 回调函数，用于自定义处理每列数据
// fieldRowNum ...int 字段行数，默认为第一行, 从1开始。
func ReadExcelRow(excelFile *excelize.File, sheetName []string, sheetIndex []int, excludeSheet []string, dataRow int, dataCol int, call func(sheetName string, tmp map[string]any) error, fieldRowNum ...int) error {
	var sheets []string
	if len(sheetName) > 0 {
		sheets = sheetName
	} else if len(sheetIndex) > 0 {
		for _, index := range sheetIndex {
			sheets = append(sheets, excelFile.GetSheetName(index))
		}
	} else {
		sheets = excelFile.GetSheetList()
	}
	var (
		errs     []string
		fieldRow int
	)
	if len(fieldRowNum) > 0 {
		fieldRow = fieldRowNum[0] - 1
	}
	var excludeSheetMap = make(map[string]bool)
	for _, sheet := range excludeSheet {
		excludeSheetMap[sheet] = true
	}
	for _, sheet := range sheets {
		if excludeSheetMap[sheet] {
			continue
		}
		rows, err := excelFile.GetRows(sheet)
		if err != nil {
			errs = append(errs, fmt.Sprintf("读取sheet%s失败: %v", sheet, err))
			continue
		}
		if len(rows) < dataRow {
			errs = append(errs, fmt.Sprintf("sheet%s数据行数小于%d，未读取到有效数据", sheet, dataRow))
			continue
		}
		fieldRows := rows[fieldRow]
		fieldLen := len(fieldRows)
		for idx, row := range rows {
			if idx < (dataRow - 1) {
				continue
			}
			if dataCol > 0 && len(row) < dataCol {
				errs = append(errs, fmt.Sprintf("sheet%s第%d行数据列数小于%d，未读取到有效数据", sheet, idx+1, dataCol))
				continue
			}
			var tmp = make(map[string]any)
			for i, cell := range row {
				if i >= fieldLen {
					break
				}
				tmp[strings.TrimSpace(fieldRows[i])] = strings.TrimSpace(cell)
			}
			if err = call(sheet, tmp); err != nil {
				errs = append(errs, fmt.Sprintf("sheet%s第%d行数据处理失败: %v", sheet, idx+1, err))
			}
		}
	}
	if len(errs) < 1 {
		return nil
	}
	return errors.New(strings.Join(errs, "\n"))
}

type HasSheetName interface {
	SheetName() string
}

// GetTplSheetFieldIndex 获取模板字段
// 读取excel里面指定sheet的数据，并提取某一列的所有字段所在的位置。
func GetTplSheetFieldIndex(file *excelize.File, model any, fieldRowNum int) (string, map[string]string, error) {
	var (
		hasSheetName HasSheetName
		sheetName    string
		ok           bool
		val          = reflect.ValueOf(model)
	)
	// 验证是否实现了wireless.HasSheetName接口
	switch val.Kind() {
	case reflect.Ptr:
		val = val.Elem()
		switch val.Kind() {
		case reflect.Slice:
			if val.Len() < 1 {
				return "", nil, errors.New("model为空")
			}
			hasSheetName, ok = val.Index(0).Interface().(HasSheetName)
		case reflect.Struct:
			hasSheetName, ok = val.Interface().(HasSheetName)
		default:
			return "", nil, errors.New("model类型错误")
		}
	case reflect.Slice:
		if val.Len() < 1 {
			return "", nil, errors.New("model为空")
		}
		hasSheetName, ok = val.Index(0).Interface().(HasSheetName)
	case reflect.Struct:
		hasSheetName, ok = val.Interface().(HasSheetName)
	case reflect.String:
		ok = true
		sheetName = val.Interface().(string)
	default:
		return "", nil, errors.New("model类型错误")

	}
	if !ok {
		return "", nil, errors.New(reflect.TypeOf(model).Name() + "未实现HasSheetName接口")
	}
	// 获取sheet内容
	if hasSheetName != nil {
		sheetName = hasSheetName.SheetName()
	}

	rows, err := file.GetRows(sheetName)
	if err != nil {
		return "", nil, err
	}
	if len(rows) < fieldRowNum {
		return "", nil, errors.New("无有效字段列")
	}
	rowValues := rows[fieldRowNum-1]
	var mapping = make(map[string]string)
	for idx, cel := range rowValues {
		celName, _ := excelize.ColumnNumberToName(idx + 1)
		mapping[cel] = celName
	}

	return sheetName, mapping, nil
}
