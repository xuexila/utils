package excelTools

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strings"
)

func CloseExcel(f *excelize.File) {
	if f != nil {
		_ = f.Close()
	}
}

func ReadExcelRow(excelFile *excelize.File, sheetName []string, sheetIndex []int, dataRow int, dataCol int, call func(tmp map[string]any) error) error {
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
	var errs []string
	for _, sheet := range sheets {
		rows, err := excelFile.GetRows(sheet)
		if err != nil {
			errs = append(errs, fmt.Sprintf("读取sheet%s失败: %v", sheet, err))
			continue
		}
		if len(rows) < dataRow {
			errs = append(errs, fmt.Sprintf("sheet%s数据行数小于%d，未读取到有效数据", sheet, dataRow))
			continue
		}
		fieldLen := len(rows[0])
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
				tmp[rows[0][i]] = cell
			}
			if err = call(tmp); err != nil {
				errs = append(errs, fmt.Sprintf("sheet%s第%d行数据处理失败: %v", sheet, idx+1, err))
			}
		}
	}
	if len(errs) < 1 {
		return nil
	}
	return errors.New(strings.Join(errs, "\n"))
}
