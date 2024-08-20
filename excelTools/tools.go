package excelTools

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
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
				tmp[fieldRows[i]] = cell
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
