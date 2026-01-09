package main

import (
	"fmt"
)

import "github.com/xuri/excelize/v2"

func ExportExcel(rowData [][]interface{}) (filename string, err error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
			return
		}
	}()
	sheet := "Sheet1"
	sw, err := f.NewStreamWriter(sheet)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	styleID, err := f.NewStyle(&excelize.Style{Font: &excelize.Font{Color: "777777"}})
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if err := sw.SetRow("A1",
		[]interface{}{
			excelize.Cell{StyleID: styleID, Value: "ID"},
			excelize.Cell{StyleID: styleID, Value: "Transaction Amount"},
			excelize.Cell{StyleID: styleID, Value: "Transaction Date"},
			excelize.Cell{StyleID: styleID, Value: "Transaction Name"},
			excelize.Cell{StyleID: styleID, Value: "Reference Number"},
		},
		excelize.RowOpts{Height: 45, Hidden: false}); err != nil {
		fmt.Println(err)
		return "", err
	}
	dataStart := 2
	for rowID, row := range rowData {
		cell, err := excelize.CoordinatesToCellName(1, dataStart+rowID)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		if err := sw.SetRow(cell, row); err != nil {
			fmt.Println(err)
			break
		}
	}
	if err := sw.Flush(); err != nil {
		fmt.Println(err)
		return "", err
	}
	fileName := "Book1.xlsx"
	if err := f.SaveAs(fileName); err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("Success Create Book")
	return fileName, nil
}
