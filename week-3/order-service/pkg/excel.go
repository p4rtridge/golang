package pkg

import (
	"fmt"
	"order_service/services/order/entity"
	"os"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

func GenerateExcel(datas *[]entity.OrdersSummarize, startDate, endDate time.Time) (string, error) {
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "summarize"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return "", err
	}
	f.DeleteSheet("Sheet1")

	err = f.MergeCell(sheetName, "A1", "E1")
	if err != nil {
		return "", err
	}

	f.SetRowHeight(sheetName, 1, 25)
	f.SetCellValue(sheetName, "A1", fmt.Sprintf("Start date: %s", startDate))

	err = f.MergeCell(sheetName, "A2", "E2")
	if err != nil {
		return "", err
	}

	f.SetRowHeight(sheetName, 2, 20)
	f.SetCellValue(sheetName, "A2", fmt.Sprintf("End date: %s", endDate))

	cells := []string{"A", "B", "C", "D", "E"}
	cellsValue := []string{"User ID", "Username", "Num of Orders", "Average Price", "Average Order Quantity"}

	cellHeaderIdx := 4
	for i := 1; i <= len(cells); i++ {
		f.SetCellValue(sheetName, fmt.Sprintf("%s%d", cells[i-1], cellHeaderIdx), cellsValue[i-1])
	}

	startRow := 5

	for _, data := range *datas {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", startRow), data.UserId)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", startRow), data.Username)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", startRow), data.NumOfOrders)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", startRow), data.SumOrderPrice)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", startRow), data.AverageOrderItemQuantity)

		startRow++
	}

	cols, err := f.GetCols(sheetName)
	if err != nil {
		return "", err
	}
	for idx, col := range cols {
		largestWidth := 0
		for _, rowCell := range col {
			cellWidth := utf8.RuneCountInString(rowCell) + 2 // + 2 for margin
			if cellWidth > largestWidth {
				largestWidth = cellWidth
			}
		}
		name, err := excelize.ColumnNumberToName(idx + 1)
		if err != nil {
			return "", err
		}
		f.SetColWidth(sheetName, name, name, float64(largestWidth))
	}

	f.SetActiveSheet(index)

	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	excelName := fmt.Sprintf("summarize-%s.xlsx", uuid.New().String())
	err = f.SaveAs(fmt.Sprintf("%s/storage/%s", pwd, excelName))
	if err != nil {
		return "", err
	}

	return excelName, nil
}
