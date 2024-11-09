package excelTools

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"reflect"
)

//
// ━━━━━━神兽出没━━━━━━
// 　　 ┏┓     ┏┓
// 　　┏┛┻━━━━━┛┻┓
// 　　┃　　　　　 ┃
// 　　┃　　━　　　┃
// 　　┃　┳┛　┗┳  ┃
// 　　┃　　　　　 ┃
// 　　┃　　┻　　　┃
// 　　┃　　　　　 ┃
// 　　┗━┓　　　┏━┛　Code is far away from bug with the animal protecting
// 　　　 ┃　　　┃    神兽保佑,代码无bug
// 　　　　┃　　　┃
// 　　　　┃　　　┗━━━┓
// 　　　　┃　　　　　　┣┓
// 　　　　┃　　　　　　┏┛
// 　　　　┗┓┓┏━┳┓┏┛
// 　　　　 ┃┫┫ ┃┫┫
// 　　　　 ┗┻┛ ┗┻┛
//
// ━━━━━━感觉萌萌哒━━━━━━
//
//
// User helay
// Date: 2024/11/9 16:19
//

// ExportToExcel 将结构体切片导出为 Excel
func ExportToExcel[T any](dst io.Writer, data []T) error {
	// 创建一个新的 Excel 文件。
	f := excelize.NewFile()

	// 获取第一个元素的类型。
	if len(data) == 0 {
		return fmt.Errorf("data slice is empty")
	}
	v := reflect.ValueOf(data[0])
	t := v.Type()

	// 创建一个工作表。
	sheetName := "Sheet1"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}
	f.SetActiveSheet(index)

	// 写入表头。
	headers := make([]string, 0)
	fieldIndices := make([]int, 0)
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		_field := field.Tag.Get("excel")
		if _field == "-" || _field == "" {
			continue
		}
		headers = append(headers, _field)
		fieldIndices = append(fieldIndices, i)
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		err := f.SetCellValue(sheetName, cell, header)
		if err != nil {
			return fmt.Errorf("set cell value failed: %v", err)
		}
	}

	// 将数据写入 Excel 文件。
	for rowIndex, item := range data {
		row := rowIndex + 2 // 跳过第一行的表头
		vItem := reflect.ValueOf(item)
		for i, fieldIndex := range fieldIndices {
			value := vItem.Field(fieldIndex).Interface()
			cell, _ := excelize.CoordinatesToCellName(i+1, row)
			err := f.SetCellValue(sheetName, cell, value)
			if err != nil {
				return fmt.Errorf("set cell value failed: %v", err)
			}
		}
	}
	return f.Write(dst)
}
