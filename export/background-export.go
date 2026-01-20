package export

import (
	"download-excel-project/repository"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/xuri/excelize/v2"
	"os"
	"runtime"
	"sync"
	"time"
)

type JobStatus struct {
	ID        string     `json:"id"`
	Status    string     `json:"status"` // "processing", "completed", "failed"
	Progress  int        `json:"progress"`
	TotalRows int32      `json:"total_rows"`
	Filename  string     `json:"filename"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	ErrorMsg  string     `json:"error_msg,omitempty"`
}

var (
	jobStore = make(map[string]*JobStatus)
	jobMutex = sync.RWMutex{}
)

func StartBackgroundExport(jobID string, totalData int32, repo repository.CustomerRepo) {
	filename := fmt.Sprintf("customer_export_%s.xlsx", jobID)
	start := time.Now()

	jobMutex.Lock()
	jobStore[jobID] = &JobStatus{
		ID:        jobID,
		Status:    "processing",
		Progress:  0,
		TotalRows: totalData,
		Filename:  filename,
		StartTime: time.Now(),
	}
	jobMutex.Unlock()

	go func() {
		var mStart, mEnd runtime.MemStats
		runtime.ReadMemStats(&mStart)

		err := exportCustomerBackground(jobID, totalData, filename, repo)

		runtime.ReadMemStats(&mEnd)
		elapsed := time.Since(start)

		jobMutex.Lock()
		job := jobStore[jobID]
		endTime := time.Now()
		job.EndTime = &endTime

		if err != nil {
			job.Status = "failed"
			job.ErrorMsg = err.Error()
		} else {
			job.Status = "completed"
			job.Progress = 100
			log.Infof("Memory Usage: Alloc = %.2f MB, TotalAlloc = %.2f MB, Sys = %.2f MB, NumGC = %v",
				float64(mEnd.Alloc-mStart.Alloc)/1024/1024,
				float64(mEnd.TotalAlloc-mStart.TotalAlloc)/1024/1024,
				float64(mEnd.Sys-mStart.Sys)/1024/1024,
				mEnd.NumGC-mStart.NumGC)
			log.Infof("Execution Time: %v", elapsed)
		}
		jobMutex.Unlock()
	}()
}

func GetJobStatus(jobID string) (*JobStatus, bool) {
	jobMutex.RLock()
	defer jobMutex.RUnlock()
	job, exists := jobStore[jobID]
	return job, exists
}

func exportCustomerBackground(jobID string, totalData int32, filename string, repo repository.CustomerRepo) error {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Customer Data"
	f.SetSheetName("Sheet1", sheet)

	sw, err := f.NewStreamWriter(sheet)
	if err != nil {
		return err
	}

	// Header
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"366092"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	sw.SetRow("A1", []interface{}{
		excelize.Cell{StyleID: headerStyle, Value: "Created Date"},
		excelize.Cell{StyleID: headerStyle, Value: "Customer Name"},
		excelize.Cell{StyleID: headerStyle, Value: "Customer Email"},
		excelize.Cell{StyleID: headerStyle, Value: "Reference Number"},
		excelize.Cell{StyleID: headerStyle, Value: "Changes"},
		excelize.Cell{StyleID: headerStyle, Value: "Old Phone"},
		excelize.Cell{StyleID: headerStyle, Value: "New Phone"},
		excelize.Cell{StyleID: headerStyle, Value: "Old Email"},
		excelize.Cell{StyleID: headerStyle, Value: "New Email"},
		excelize.Cell{StyleID: headerStyle, Value: "Old Occupation"},
		excelize.Cell{StyleID: headerStyle, Value: "New Occupation"},
		excelize.Cell{StyleID: headerStyle, Value: "Old Job"},
		excelize.Cell{StyleID: headerStyle, Value: "New Job"},
		excelize.Cell{StyleID: headerStyle, Value: "Old Company"},
		excelize.Cell{StyleID: headerStyle, Value: "New Company"},
		excelize.Cell{StyleID: headerStyle, Value: "Old Company Address"},
		excelize.Cell{StyleID: headerStyle, Value: "New Company Address"},
		excelize.Cell{StyleID: headerStyle, Value: "Old Company Type"},
		excelize.Cell{StyleID: headerStyle, Value: "New Company Type"},
		excelize.Cell{StyleID: headerStyle, Value: "Old Income"},
		excelize.Cell{StyleID: headerStyle, Value: "New Income"},
		excelize.Cell{StyleID: headerStyle, Value: "Old Emergency Phone"},
		excelize.Cell{StyleID: headerStyle, Value: "New Emergency Phone"},
		excelize.Cell{StyleID: headerStyle, Value: "Old Emergency Contact"},
		excelize.Cell{StyleID: headerStyle, Value: "New Emergency Contact"},
		excelize.Cell{StyleID: headerStyle, Value: "Old Address"},
		excelize.Cell{StyleID: headerStyle, Value: "New Address"},
		excelize.Cell{StyleID: headerStyle, Value: "Old Kecamatan"},
		excelize.Cell{StyleID: headerStyle, Value: "New Kecamatan"},
		excelize.Cell{StyleID: headerStyle, Value: "Old Kelurahan"},
		excelize.Cell{StyleID: headerStyle, Value: "New Kelurahan"},
		excelize.Cell{StyleID: headerStyle, Value: "Old RT"},
		excelize.Cell{StyleID: headerStyle, Value: "New RT"},
		excelize.Cell{StyleID: headerStyle, Value: "Old RW"},
		excelize.Cell{StyleID: headerStyle, Value: "New RW"},
		excelize.Cell{StyleID: headerStyle, Value: "Old City"},
		excelize.Cell{StyleID: headerStyle, Value: "New City"},
		excelize.Cell{StyleID: headerStyle, Value: "Old Province"},
		excelize.Cell{StyleID: headerStyle, Value: "New Province"},
		excelize.Cell{StyleID: headerStyle, Value: "Old Postal Code"},
		excelize.Cell{StyleID: headerStyle, Value: "New Postal Code"},
	})

	dataChan, errChan := repo.GenerateUpdateDataWithCursor(totalData)
	rowNum := 2
	processed := 0

	for {
		select {
		case row, ok := <-dataChan:
			if !ok {
				sw.Flush()
				if err := f.SaveAs(filename); err != nil {
					return err
				}
				return nil
			}

			cell, _ := excelize.CoordinatesToCellName(1, rowNum)
			sw.SetRow(cell, row)

			rowNum++
			processed++

			// Update progress every 10k records
			if processed%10000 == 0 {
				progress := int((float64(processed) / float64(totalData)) * 100)
				jobMutex.Lock()
				if job, exists := jobStore[jobID]; exists {
					job.Progress = progress
				}
				jobMutex.Unlock()
			}

		case err := <-errChan:
			if err != nil {
				return err
			}
		case <-time.After(45 * time.Minute):
			return fmt.Errorf("export timeout")
		}
	}
}

func CleanupJob(jobID string) {
	jobMutex.Lock()
	defer jobMutex.Unlock()

	if job, exists := jobStore[jobID]; exists {
		// Remove file
		os.Remove(job.Filename)
		// Remove from store
		delete(jobStore, jobID)
	}
}
