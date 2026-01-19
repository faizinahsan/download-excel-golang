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
