package main

import (
	"download-excel-project/config"
	"download-excel-project/export"
	"download-excel-project/repository"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"os"
	"runtime"
	"strconv"
	"time"
)

type Handler struct {
	transactionRepo repository.TransactionRepo
	customerRepo    repository.CustomerRepo
}

func main() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.SetConfigName("env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}
	app := fiber.New(fiber.Config{
		Prefork:       false, // Disabled for compatibility
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Test App v1.0.1",
	})
	db, err := config.ConnectToDB()
	if err != nil {
		panic(err)
	}
	transactionRepo := repository.NewTransactionRepo(db)
	customerRepo := repository.NewCustomerRepo(db)
	handler := NewServiceHandler(transactionRepo, customerRepo)

	app.Get("/hello", handler.HelloWorld)
	transaction := app.Group("/transactions")
	transaction.Get("/generate-file", handler.GetTransactions)
	transaction.Get("/download", handler.DownloadTransactionFile)

	customer := app.Group("/customers")
	customer.Get("/generate-update-data", handler.GenerateAuditTrailData)
	customer.Get("/download-update-data", handler.DownloadAuditTrailFile)

	// New streaming Excel export endpoints
	customer.Post("/export/excel/streaming", handler.ExportCustomerExcelStreaming)
	customer.Get("/export/excel/streaming", handler.ExportCustomerExcelStreamingQuery)
	customer.Get("/download/:filename", handler.DownloadCustomerExcelFile)

	// Background export endpoints
	customer.Post("/export/background", handler.StartBackgroundExport)
	customer.Get("/export/status/:jobId", handler.GetExportStatus)
	customer.Get("/export/download/:jobId", handler.DownloadBackgroundExport)

	port := viper.GetString("server.port")
	if port == "" {
		log.Fatalf("Fatal error: server.port is not set in config file")
	}
	log.Infof("Attempting to start server on port: %s", port)

	err = app.Listen(":" + port)
	if err != nil {
		log.Errorf("Failed to start server on port %s: %T: %v | %s", port, err, err, err.Error())
		os.Exit(1)
	}
}

func NewServiceHandler(
	transactionRepo repository.TransactionRepo,
	customerRepo repository.CustomerRepo,
) *Handler {
	return &Handler{
		transactionRepo: transactionRepo,
		customerRepo:    customerRepo,
	}
}

func (h *Handler) HelloWorld(c *fiber.Ctx) error {
	log.Info("Request received for HelloWorld endpoint")
	return c.SendString("Hello, World!")
}

func (h *Handler) GetTransactions(c *fiber.Ctx) error {
	start := time.Now()
	var mStart, mEnd runtime.MemStats
	runtime.ReadMemStats(&mStart)
	log.Info("Received request to generate transactions")
	go func() {

	}()

	totalData := 0
	if totalDataReq := c.Query("total_data"); totalDataReq != "" {
		log.Infof("Query parameter total_data: %s", totalDataReq)
		totalData, _ = strconv.Atoi(totalDataReq)
	} else {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	transactions, err := h.transactionRepo.GetTransactionToArray(int32(totalData))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	log.Infof("Retrieved %d transactions", len(transactions))
	filename, err := export.ExportExcel(transactions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	zipFile := filename + ".zip"
	err = export.AddToZip(filename, zipFile)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	log.Infof("Generated file: %s", zipFile)

	runtime.ReadMemStats(&mEnd)
	elapsed := time.Since(start)
	log.Infof("Memory Usage: Alloc = %.2f MB, TotalAlloc = %.2f MB, Sys = %.2f MB, NumGC = %v", float64(mEnd.Alloc-mStart.Alloc)/1024/1024, float64(mEnd.TotalAlloc-mStart.TotalAlloc)/1024/1024, float64(mEnd.Sys-mStart.Sys)/1024/1024, mEnd.NumGC-mStart.NumGC)
	defer func() {
		log.Infof("Execution Time: %v", elapsed)
	}()
	return c.Status(200).JSON(fiber.Map{
		"message": "File generated successfully",
		"file":    zipFile,
	})
}

func (h *Handler) DownloadTransactionFile(c *fiber.Ctx) error {
	start := time.Now()
	var mStart, mEnd runtime.MemStats
	runtime.ReadMemStats(&mStart)
	log.Info("Received request to download transactions")
	filename := c.Query("filename")
	defer func() {
		runtime.ReadMemStats(&mEnd)
		elapsed := time.Since(start)
		log.Infof("Memory Usage: Alloc = %.2f MB, TotalAlloc = %.2f MB, Sys = %.2f MB, NumGC = %v", float64(mEnd.Alloc-mStart.Alloc)/1024/1024, float64(mEnd.TotalAlloc-mStart.TotalAlloc)/1024/1024, float64(mEnd.Sys-mStart.Sys)/1024/1024, mEnd.NumGC-mStart.NumGC)
		err := os.Remove("./" + filename)
		if err != nil {
			log.Errorf("Failed to remove file %s: %v", filename, err)
		} else {
			log.Infof("Successfully removed file %s", filename)
		}
		log.Infof("Execution Time: %v", elapsed)
	}()
	return c.Download("./" + filename)
}

func (h *Handler) GenerateAuditTrailData(c *fiber.Ctx) error {
	start := time.Now()
	var mStart, mEnd runtime.MemStats
	runtime.ReadMemStats(&mStart)
	log.Info("Received request to generate customer update data")
	fmt.Println("=== Method 1: Batch Processing ===")

	//totalData := 0
	//if totalDataReq := c.Query("total_data"); totalDataReq != "" {
	//	log.Infof("Query parameter total_data: %s", totalDataReq)
	//	totalData, _ = strconv.Atoi(totalDataReq)
	//} else {
	//	return c.SendStatus(fiber.StatusBadRequest)
	//}
	filename := "AuditTrailProfile1.xlsx"
	zipFile := filename + ".zip"
	// Method 1: Batch Processing - Sangat direkomendasikan untuk 1 juta data
	totalDataReq, err := strconv.Atoi(c.Query("total_data"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	totalData := int32(totalDataReq)

	go func() {
		var allData [][]interface{}
		var batchSize int32 = 50000 // Process 50k records per batch
		var offset int32 = 0

		for offset < totalData {
			currentBatch := batchSize
			if offset+batchSize > totalData {
				currentBatch = totalData - offset
			}

			log.Infof("Processing batch: offset=%d, size=%d\n", offset, currentBatch)

			batch, err := h.customerRepo.GenerateUpdateDataOptimized(currentBatch, offset)
			if err != nil {
				log.Fatal(err)
			}

			allData = append(allData, batch...)
			offset += batchSize
		}
		filename, err := export.ExportAuditTrailToExcelV2(allData)
		if err != nil {
			log.Errorf("Error exporting audit trail to Excel: %v", err)
			return
		}
		err = export.AddToZip(filename, zipFile)
		if err != nil {
			log.Errorf("Error creating zip file: %v", err)
			return
		}
		log.Infof("Generated file: %s", zipFile)

		runtime.ReadMemStats(&mEnd)
		elapsed := time.Since(start)
		log.Infof("Memory Usage: Alloc = %.2f MB, TotalAlloc = %.2f MB, Sys = %.2f MB, NumGC = %v", float64(mEnd.Alloc-mStart.Alloc)/1024/1024, float64(mEnd.TotalAlloc-mStart.TotalAlloc)/1024/1024, float64(mEnd.Sys-mStart.Sys)/1024/1024, mEnd.NumGC-mStart.NumGC)
		defer func() {
			log.Infof("Execution Time: %v", elapsed)
		}()
	}()
	return c.Status(200).JSON(fiber.Map{
		"message": "Customer update data file generated successfully",
		"file":    zipFile,
	})
}

func (h *Handler) DownloadAuditTrailFile(c *fiber.Ctx) error {
	start := time.Now()
	var mStart, mEnd runtime.MemStats
	runtime.ReadMemStats(&mStart)
	log.Info("Received request to download transactions")
	filename := c.Query("filename")
	defer func() {
		runtime.ReadMemStats(&mEnd)
		elapsed := time.Since(start)
		log.Infof("Memory Usage: Alloc = %.2f MB, TotalAlloc = %.2f MB, Sys = %.2f MB, NumGC = %v", float64(mEnd.Alloc-mStart.Alloc)/1024/1024, float64(mEnd.TotalAlloc-mStart.TotalAlloc)/1024/1024, float64(mEnd.Sys-mStart.Sys)/1024/1024, mEnd.NumGC-mStart.NumGC)
		err := os.Remove("./" + filename)
		if err != nil {
			log.Errorf("Failed to remove file %s: %v", filename, err)
		} else {
			log.Infof("Successfully removed file %s", filename)
		}
		log.Infof("Execution Time: %v", elapsed)
	}()
	return c.Download("./" + filename)
}

// ExportCustomerExcelStreaming - Streaming Excel export dengan POST request body
func (h *Handler) ExportCustomerExcelStreaming(c *fiber.Ctx) error {
	start := time.Now()
	var mStart, mEnd runtime.MemStats
	runtime.ReadMemStats(&mStart)
	log.Info("Received request for streaming Excel export (POST)")

	// Parse request body
	var req struct {
		TotalData int32  `json:"total_data"`
		Filename  string `json:"filename"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate request
	if req.TotalData <= 0 || req.TotalData > 5000000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "total_data must be between 1 and 5,000,000",
		})
	}

	// Generate filename if not provided
	filename := req.Filename
	if filename == "" {
		filename = fmt.Sprintf("customer_updates_streaming_%s.xlsx", time.Now().Format("20060102_150405"))
	}

	// Ensure .xlsx extension
	if len(filename) < 5 || filename[len(filename)-5:] != ".xlsx" {
		filename += ".xlsx"
	}

	// Start streaming export
	resultFilename, err := export.ExportCustomerToExcelWithStreaming(h.customerRepo, req.TotalData, filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to generate Excel file: %v", err),
		})
	}

	runtime.ReadMemStats(&mEnd)
	elapsed := time.Since(start)
	log.Infof("Memory Usage: Alloc = %.2f MB, TotalAlloc = %.2f MB, Sys = %.2f MB, NumGC = %v",
		float64(mEnd.Alloc-mStart.Alloc)/1024/1024,
		float64(mEnd.TotalAlloc-mStart.TotalAlloc)/1024/1024,
		float64(mEnd.Sys-mStart.Sys)/1024/1024,
		mEnd.NumGC-mStart.NumGC)
	log.Infof("Execution Time: %v", elapsed)

	return c.JSON(fiber.Map{
		"success":    true,
		"message":    "Excel file generated successfully with streaming",
		"filename":   resultFilename,
		"total_rows": req.TotalData,
		"duration":   elapsed.String(),
	})
}

// ExportCustomerExcelStreamingQuery - Streaming Excel export dengan query parameters
func (h *Handler) ExportCustomerExcelStreamingQuery(c *fiber.Ctx) error {
	start := time.Now()
	var mStart, mEnd runtime.MemStats
	runtime.ReadMemStats(&mStart)
	log.Info("Received request for streaming Excel export (GET)")

	// Parse query parameters
	totalDataStr := c.Query("total_data")
	if totalDataStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "total_data query parameter is required",
		})
	}

	totalData, err := strconv.ParseInt(totalDataStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid total_data parameter",
		})
	}

	// Validate limits
	if totalData <= 0 || totalData > 5000000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "total_data must be between 1 and 5,000,000",
		})
	}

	// Generate filename
	filename := c.Query("filename")
	if filename == "" {
		filename = fmt.Sprintf("customer_updates_streaming_%s.xlsx", time.Now().Format("20060102_150405"))
	}

	// Ensure .xlsx extension
	if len(filename) < 5 || filename[len(filename)-5:] != ".xlsx" {
		filename += ".xlsx"
	}

	// Start streaming export
	resultFilename, err := export.ExportCustomerToExcelWithStreaming(h.customerRepo, int32(totalData), filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to generate Excel file: %v", err),
		})
	}

	runtime.ReadMemStats(&mEnd)
	elapsed := time.Since(start)
	log.Infof("Memory Usage: Alloc = %.2f MB, TotalAlloc = %.2f MB, Sys = %.2f MB, NumGC = %v",
		float64(mEnd.Alloc-mStart.Alloc)/1024/1024,
		float64(mEnd.TotalAlloc-mStart.TotalAlloc)/1024/1024,
		float64(mEnd.Sys-mStart.Sys)/1024/1024,
		mEnd.NumGC-mStart.NumGC)
	log.Infof("Execution Time: %v", elapsed)

	return c.JSON(fiber.Map{
		"success":    true,
		"message":    "Excel file generated successfully with streaming (query params)",
		"filename":   resultFilename,
		"total_rows": int32(totalData),
		"duration":   elapsed.String(),
	})
}

// DownloadCustomerExcelFile - Download file Excel yang sudah dibuat
func (h *Handler) DownloadCustomerExcelFile(c *fiber.Ctx) error {
	start := time.Now()
	var mStart, mEnd runtime.MemStats
	runtime.ReadMemStats(&mStart)
	log.Info("Received request to download customer Excel file")

	filename := c.Params("filename")
	if filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "filename parameter is required",
		})
	}

	defer func() {
		runtime.ReadMemStats(&mEnd)
		elapsed := time.Since(start)
		log.Infof("Memory Usage: Alloc = %.2f MB, TotalAlloc = %.2f MB, Sys = %.2f MB, NumGC = %v",
			float64(mEnd.Alloc-mStart.Alloc)/1024/1024,
			float64(mEnd.TotalAlloc-mStart.TotalAlloc)/1024/1024,
			float64(mEnd.Sys-mStart.Sys)/1024/1024,
			mEnd.NumGC-mStart.NumGC)
		log.Infof("Download Execution Time: %v", elapsed)

		// Optionally remove file after download (uncomment if desired)
		// err := os.Remove("./" + filename)
		// if err != nil {
		//     log.Errorf("Failed to remove file %s: %v", filename, err)
		// } else {
		//     log.Infof("Successfully removed file %s", filename)
		// }
	}()

	// Set headers for Excel file download
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	return c.SendFile("./" + filename)
}

// StartBackgroundExport - Start background Excel export
func (h *Handler) StartBackgroundExport(c *fiber.Ctx) error {
	var req struct {
		TotalData int32 `json:"total_data"`
	}

	if err := c.BodyParser(&req); err != nil {
		log.Errorf("Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if req.TotalData <= 0 || req.TotalData > 5000000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "total_data must be between 1 and 5,000,000",
		})
	}

	jobID := fmt.Sprintf("%d", time.Now().UnixNano())
	export.StartBackgroundExport(jobID, req.TotalData, h.customerRepo)

	return c.JSON(fiber.Map{
		"success": true,
		"job_id":  jobID,
		"message": "Export started in background",
	})
}

// GetExportStatus - Get background export status
func (h *Handler) GetExportStatus(c *fiber.Ctx) error {
	jobID := c.Params("jobId")
	if jobID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "job_id is required",
		})
	}

	status, exists := export.GetJobStatus(jobID)
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Job not found",
		})
	}

	return c.JSON(status)
}

// DownloadBackgroundExport - Download completed background export
func (h *Handler) DownloadBackgroundExport(c *fiber.Ctx) error {
	jobID := c.Params("jobId")
	if jobID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "job_id is required",
		})
	}

	status, exists := export.GetJobStatus(jobID)
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Job not found",
		})
	}

	if status.Status != "completed" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Export not ready. Status: %s", status.Status),
		})
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", status.Filename))

	// Cleanup after download
	defer export.CleanupJob(jobID)

	return c.SendFile("./" + status.Filename)
}
