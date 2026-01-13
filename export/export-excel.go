package export

import (
	"download-excel-project/repository"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
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

func GenerateCustomerUpdateDataExcel(rowData [][]interface{}) (filename string, err error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Errorf(err.Error())
			return
		}
	}()
	sheet := "Sheet1"
	sw, err := f.NewStreamWriter(sheet)
	if err != nil {
		log.Errorf(err.Error())
		return "", err
	}
	styleID, err := f.NewStyle(&excelize.Style{Font: &excelize.Font{Color: "777777"}})
	if err != nil {
		log.Errorf(err.Error())

		return "", err
	}
	if err := sw.SetRow("A1",
		[]interface{}{
			excelize.Cell{StyleID: styleID, Value: "Created"},
			excelize.Cell{StyleID: styleID, Value: "Customer Name"},
			excelize.Cell{StyleID: styleID, Value: "Customer Email"},
			excelize.Cell{StyleID: styleID, Value: "Reference Number"},
			excelize.Cell{StyleID: styleID, Value: "List Data Changes"},
		},
		excelize.RowOpts{Height: 45, Hidden: false}); err != nil {
		log.Errorf(err.Error())
		return "", err
	}
	dataStart := 3
	for rowID, row := range rowData {
		cell, err := excelize.CoordinatesToCellName(1, dataStart+rowID)
		if err != nil {
			log.Errorf(err.Error())
			return "", err
		}
		if err := sw.SetRow(cell, row); err != nil {
			log.Errorf(err.Error())
			break
		}
	}
	if err := sw.Flush(); err != nil {
		log.Errorf(err.Error())
		return "", err
	}
	fileName := "AuditTrailProfile1.xlsx"
	if err := f.SaveAs(fileName); err != nil {
		log.Errorf(err.Error())
		return "", err
	}
	log.Info("Success Create Book")
	return fileName, nil
}

func ExportToExcelWithStreaming(repo repository.CustomerRepo, totalData int32) (string, error) {
	dataChan, errChan := repo.GenerateUpdateDataWithCursor(totalData)
	fileName := "AuditTrailProfile1.xlsx"
	// Setup Excel file (contoh dengan excelize)
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Errorf(err.Error())
			return
		}
	}()
	sheet := "Sheet1"
	sw, err := f.NewStreamWriter(sheet)
	if err != nil {
		log.Errorf(err.Error())
		return "", err
	}
	styleID, err := f.NewStyle(&excelize.Style{Font: &excelize.Font{Color: "777777"}})
	if err != nil {
		log.Errorf(err.Error())

		return "", err
	}
	if err := sw.SetRow("A1",
		[]interface{}{
			excelize.Cell{StyleID: styleID, Value: "Created"},
			excelize.Cell{StyleID: styleID, Value: "Customer Name"},
			excelize.Cell{StyleID: styleID, Value: "Customer Email"},
			excelize.Cell{StyleID: styleID, Value: "Reference Number"},
			excelize.Cell{StyleID: styleID, Value: "List Data Changes"},
		},
		excelize.RowOpts{Height: 45, Hidden: false}); err != nil {
		log.Errorf(err.Error())
		return "", err
	}

	var rowNum = 2 // Start from row 2 (row 1 untuk header)

	for {
		select {
		case _, ok := <-dataChan:
			if !ok {
				// Channel closed, processing completed
				fmt.Printf("Excel export selesai: total %d rows\n", rowNum-2)
				return fileName, nil
			}

			// Write to Excel
			// for colIndex, value := range row {
			//     cell := fmt.Sprintf("%s%d", getColumnName(colIndex), rowNum)
			//     f.SetCellValue("Sheet1", cell, value)
			// }

			rowNum++

			if rowNum%10000 == 0 {
				fmt.Printf("Written %d rows to Excel...\n", rowNum-2)
			}

		case err := <-errChan:
			if err != nil {
				return "", fmt.Errorf("error during streaming: %v", err)
			}
		}
	}
}
