package excel

import (
	"fmt"
	"io"
	"strings"

	"github.com/xuri/excelize/v2"
)

const defaultSheet = "Sheet1"

// WriteSheet 按表头与行数据生成 Excel 文件。headers 为第一行，rows 为数据行，写入默认 Sheet1。
func WriteSheet(headers []string, rows [][]string) (*excelize.File, error) {
	f := excelize.NewFile()
	sheet := defaultSheet
	for col, h := range headers {
		cell, err := excelize.CoordinatesToCellName(col+1, 1)
		if err != nil {
			return nil, err
		}
		if err := f.SetCellValue(sheet, cell, h); err != nil {
			return nil, err
		}
	}
	for r, row := range rows {
		for c, val := range row {
			if c >= len(headers) {
				break
			}
			cell, err := excelize.CoordinatesToCellName(c+1, r+2)
			if err != nil {
				return nil, err
			}
			if err := f.SetCellValue(sheet, cell, val); err != nil {
				return nil, err
			}
		}
	}
	return f, nil
}

// ReadSheet 读取指定 Sheet 的表头（第一行）与数据行，数据行以 map 形式返回，键为表头名。
func ReadSheet(f *excelize.File, sheetName string) (headers []string, rows []map[string]string, err error) {
	if sheetName == "" {
		sheetName = defaultSheet
	}
	rawRows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, nil, err
	}
	if len(rawRows) == 0 {
		return nil, []map[string]string{}, nil
	}
	headers = trimStrings(rawRows[0])
	if len(headers) == 0 {
		return nil, nil, fmt.Errorf("excel: sheet %s has empty header row", sheetName)
	}
	rows = make([]map[string]string, 0, len(rawRows)-1)
	for i := 1; i < len(rawRows); i++ {
		raw := rawRows[i]
		row := make(map[string]string, len(headers))
		for j, h := range headers {
			var v string
			if j < len(raw) {
				v = strings.TrimSpace(raw[j])
			}
			row[h] = v
		}
		rows = append(rows, row)
	}
	return headers, rows, nil
}

// StyleHeaderRow 对表头行（第 1 行）设置加粗、浅灰底、居中，并为各列设置宽度。colWidths 可为 nil，则全部用 defaultColWidth。
func StyleHeaderRow(f *excelize.File, sheet string, numCols int, colWidths []float64) error {
	const defaultColWidth = 14.0
	if sheet == "" {
		sheet = defaultSheet
	}
	// 表头样式：加粗、浅灰底、水平居中、细边框
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#E8E8E8"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "#CCCCCC", Style: 1},
			{Type: "top", Color: "#CCCCCC", Style: 1},
			{Type: "bottom", Color: "#CCCCCC", Style: 1},
			{Type: "right", Color: "#CCCCCC", Style: 1},
		},
	})
	if err != nil {
		return err
	}
	startCell, _ := excelize.CoordinatesToCellName(1, 1)
	endCell, err := excelize.CoordinatesToCellName(numCols, 1)
	if err != nil {
		return err
	}
	if err := f.SetCellStyle(sheet, startCell, endCell, headerStyle); err != nil {
		return err
	}
	// 列宽
	for i := 1; i <= numCols; i++ {
		w := defaultColWidth
		if colWidths != nil && i-1 < len(colWidths) {
			w = colWidths[i-1]
		}
		if w <= 0 {
			w = defaultColWidth
		}
		colName, err := excelize.ColumnNumberToName(i)
		if err != nil {
			return err
		}
		if err := f.SetColWidth(sheet, colName, colName, w); err != nil {
			return err
		}
	}
	return nil
}

// ReadSheetFromReader 从 io.Reader 读取 xlsx 并解析默认 Sheet1 的表头与数据行。
func ReadSheetFromReader(r io.Reader) (headers []string, rows []map[string]string, err error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	return ReadSheet(f, defaultSheet)
}

func trimStrings(ss []string) []string {
	out := make([]string, len(ss))
	for i, s := range ss {
		out[i] = strings.TrimSpace(s)
	}
	return out
}
