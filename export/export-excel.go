package export

import (
	"download-excel-project/constants"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"time"
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

func GenerateCustomerUpdateDataExcel(rowData [][]interface{}, fileName string) (filename string, err error) {
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
	if err := f.SaveAs(fileName); err != nil {
		log.Errorf(err.Error())
		return "", err
	}
	log.Info("Success Create Book")
	return fileName, nil
}

func ExportAuditTrailToExcelV2(rowData [][]interface{}) (string, error) {
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
			excelize.Cell{StyleID: styleID, Value: "Old Phone Number"},
			excelize.Cell{StyleID: styleID, Value: "New Phone Number"},
			excelize.Cell{StyleID: styleID, Value: "Old Email"},
			excelize.Cell{StyleID: styleID, Value: "New Email"},
			excelize.Cell{StyleID: styleID, Value: "Old Address Domicile"},
			excelize.Cell{StyleID: styleID, Value: "New Address Domicile"},
			excelize.Cell{StyleID: styleID, Value: "Old Occupation"},
			excelize.Cell{StyleID: styleID, Value: "New Occupation"},
			excelize.Cell{StyleID: styleID, Value: "Old Job Title"},
			excelize.Cell{StyleID: styleID, Value: "New Job Title"},
			excelize.Cell{StyleID: styleID, Value: "Old Company"},
			excelize.Cell{StyleID: styleID, Value: "New Company"},
			excelize.Cell{StyleID: styleID, Value: "Old Company Address"},
			excelize.Cell{StyleID: styleID, Value: "New Company Address"},
			excelize.Cell{StyleID: styleID, Value: "Old Company Type"},
			excelize.Cell{StyleID: styleID, Value: "New Company Type"},
			excelize.Cell{StyleID: styleID, Value: "Old Income Range"},
			excelize.Cell{StyleID: styleID, Value: "New Income Range"},
			excelize.Cell{StyleID: styleID, Value: "Old Emergency Contact"},
			excelize.Cell{StyleID: styleID, Value: "New Emergency Contact"},
			excelize.Cell{StyleID: styleID, Value: "Old Emergency Contact Name"},
			excelize.Cell{StyleID: styleID, Value: "New Emergency Contact Name"},
		},
		excelize.RowOpts{Height: 45, Hidden: false}); err != nil {
		log.Errorf(err.Error())
		return "", err
	}
	dataStart := 2
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
	if err := f.SaveAs(fileName); err != nil {
		log.Errorf(err.Error())
		return "", err
	}
	log.Info("Success Create Book")
	return fileName, nil
}
func ExportAuditTrailToExcel(rows [][]interface{}, totalData int32) (string, error) {
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
	// =========================
	// TITLE
	// =========================
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	if err := f.SetCellValue(sheet, "A1", "Audit Trail Customer"); err != nil {
		return "", err
	}
	if err := f.SetCellValue(sheet, "A2", "Export Date : "+time.Now().Format("02/01/2006")); err != nil {
		return "", err
	}
	if err := f.SetCellStyle(sheet, "A1", "A2", titleStyle); err != nil {
		return "", err
	}
	// =========================
	// HEADER ROW INDEX
	// =========================
	headerTop := 4
	headerSub := 5
	dataStart := 6

	// =========================
	// MERGE HEADER
	// =========================
	// Left headers (vertical merge)
	leftCols := []string{"A", "B", "C", "D", "E", "F"}
	for _, col := range leftCols {
		if err := f.MergeCell(sheet, col+fmt.Sprint(headerTop), col+fmt.Sprint(headerSub)); err != nil {
			return "", err
		}
	}

	// Group headers
	if err := f.MergeCell(sheet, "G4", "H4"); err != nil {
		return "", err
	}
	if err := f.MergeCell(sheet, "I4", "J4"); err != nil {
		return "", err
	}
	if err := f.MergeCell(sheet, "K4", "L4"); err != nil {
		return "", err
	}
	if err := f.MergeCell(sheet, "M4", "N4"); err != nil {
		return "", err
	}
	if err := f.MergeCell(sheet, "O4", "P4"); err != nil {
		return "", err
	}
	if err := f.MergeCell(sheet, "Q4", "R4"); err != nil {
		return "", err
	}
	if err := f.MergeCell(sheet, "S4", "T4"); err != nil {
		return "", err
	}
	if err := f.MergeCell(sheet, "U4", "V4"); err != nil {
		return "", err
	}
	if err := f.MergeCell(sheet, "W4", "X4"); err != nil {
		return "", err
	}
	if err := f.MergeCell(sheet, "Y4", "Z4"); err != nil {
		return "", err
	}
	if err := f.MergeCell(sheet, "AA4", "AB4"); err != nil {
		return "", err
	}

	if err := sw.SetRow("A4",
		[]interface{}{
			excelize.Cell{StyleID: titleStyle, Value: "No"},
			excelize.Cell{StyleID: titleStyle, Value: "Created"},
			excelize.Cell{StyleID: titleStyle, Value: "Customer Name"},
			excelize.Cell{StyleID: titleStyle, Value: "Customer Email"},
			excelize.Cell{StyleID: titleStyle, Value: "Reference Number"},
			excelize.Cell{StyleID: titleStyle, Value: "List Data Changes"},
		},
		excelize.RowOpts{Height: 45, Hidden: false}); err != nil {
		log.Errorf(err.Error())
		return "", err
	}

	// Group titles
	if err := f.SetCellValue(sheet, "G4", constants.CellularNoField); err != nil {
		return "", err
	}
	if err := f.SetCellValue(sheet, "I4", constants.EmailField); err != nil {
		return "", err
	}
	if err := f.SetCellValue(sheet, "K4", constants.AddressDomicileField); err != nil {
		return "", err
	}
	if err := f.SetCellValue(sheet, "M4", constants.OccupationField); err != nil {
		return "", err
	}
	if err := f.SetCellValue(sheet, "O4", constants.JobTitleField); err != nil {
		return "", err
	}
	if err := f.SetCellValue(sheet, "Q4", constants.CompanyField); err != nil {
		return "", err
	}
	if err := f.SetCellValue(sheet, "S4", constants.CompanyAddressField); err != nil {
		return "", err
	}
	if err := f.SetCellValue(sheet, "U4", constants.CompanyTypeField); err != nil {
		return "", err
	}
	if err := f.SetCellValue(sheet, "W4", constants.IncomeRangeField); err != nil {
		return "", err
	}
	if err := f.SetCellValue(sheet, "Y4", constants.EmergencyContactField); err != nil {
		return "", err
	}
	if err := f.SetCellValue(sheet, "AA4", constants.EmergencyContactNameField); err != nil {
		return "", err
	}

	// Sub headers
	subHeaders := []struct {
		cell  string
		value string
	}{
		{"G5", "Data Lama"}, {"H5", "Data Baru"},
		{"I5", "Data Lama"}, {"J5", "Data Baru"},
		{"K5", "Data Lama"}, {"L5", "Data Baru"},
		{"M5", "Data Lama"}, {"N5", "Data Baru"},
		{"O5", "Data Lama"}, {"P5", "Data Baru"},
		{"Q5", "Data Lama"}, {"R5", "Data Baru"},
		{"S5", "Data Lama"}, {"T5", "Data Baru"},
		{"U5", "Data Lama"}, {"V5", "Data Baru"},
		{"W5", "Data Lama"}, {"X5", "Data Baru"},
		{"Y5", "Data Lama"}, {"Z5", "Data Baru"},
		{"AA5", "Data Lama"}, {"AB5", "Data Baru"},
	}
	for _, h := range subHeaders {
		if err := f.SetCellValue(sheet, h.cell, h.value); err != nil {
			return "", err
		}
	}
	// =========================
	// STYLES
	// =========================
	border := []excelize.Border{
		{Type: "left", Style: 1},
		{Type: "right", Style: 1},
		{Type: "top", Style: 1},
		{Type: "bottom", Style: 1},
	}

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#3F00B8"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: border,
	})

	dataStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical: "top",
			WrapText: true,
		},
		Border: border,
	})

	// Apply header styles
	if err := f.SetCellStyle(sheet, "A4", "AB5", headerStyle); err != nil {
		return "", err
	}

	// =========================
	// COLUMN WIDTH
	// =========================
	if err := f.SetColWidth(sheet, "A", "A", 20); err != nil {
		return "", err
	}
	if err := f.SetColWidth(sheet, "B", "B", 18); err != nil {
		return "", err
	}
	if err := f.SetColWidth(sheet, "C", "C", 25); err != nil {
		return "", err
	}
	if err := f.SetColWidth(sheet, "D", "D", 15); err != nil {
		return "", err
	}
	if err := f.SetColWidth(sheet, "E", "E", 20); err != nil {
		return "", err
	}
	if err := f.SetColWidth(sheet, "F", "F", 30); err != nil {
		return "", err
	}
	if err := f.SetColWidth(sheet, "G", "AB", 18); err != nil {
		return "", err
	}
	// =========================
	// DATA
	// =========================
	for r, row := range rows {
		for c, val := range row {
			cell, _ := excelize.CoordinatesToCellName(c+1, dataStart+r) // c+1 karena kolom Excel mulai dari 1
			if err := f.SetCellValue(sheet, cell, val); err != nil {
				return "", err
			}
			if err := f.SetCellStyle(sheet, cell, cell, dataStyle); err != nil {
				return "", err
			}
		}
		// style empty data-lama/baru cells
		if err := f.SetCellStyle(sheet,
			fmt.Sprintf("G%d", dataStart+r),
			fmt.Sprintf("AB%d", dataStart+r),
			dataStyle,
		); err != nil {
			return "", err
		}
	}

	// =========================
	// SAVE FILE
	// =========================
	if err := f.SaveAs(fileName); err != nil {
		return "", err
	}
	return fileName, nil

}
