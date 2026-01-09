package main

import (
	"download-excel-project/config"
	"download-excel-project/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"os"
	"runtime"
	"time"
)

type Handler struct {
	transactionRepo repository.TransactionRepo
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
		Prefork:       true,
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
	handler := NewServiceHandler(transactionRepo)

	app.Get("/hello", handler.HelloWorld)
	transaction := app.Group("/transactions")
	transaction.Get("/", handler.GetTransactions)

	err = app.Listen(":" + viper.GetString("server.port"))
	if err != nil {
		log.Errorf("Failed to start server: %v", err)
	}
}

func NewServiceHandler(transactionRepo repository.TransactionRepo) *Handler {
	return &Handler{
		transactionRepo: transactionRepo,
	}
}

func (h *Handler) HelloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func (h *Handler) GetTransactions(c *fiber.Ctx) error {
	start := time.Now()
	var mStart, mEnd runtime.MemStats
	runtime.ReadMemStats(&mStart)

	transactions, err := h.transactionRepo.GetTransactionToArray()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	log.Infof("Retrieved %d transactions", len(transactions))
	filename, err := ExportExcel(transactions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	runtime.ReadMemStats(&mEnd)
	elapsed := time.Since(start)
	log.Infof("Memory Usage: Alloc = %.2f MB, TotalAlloc = %.2f MB, Sys = %.2f MB, NumGC = %v", float64(mEnd.Alloc-mStart.Alloc)/1024/1024, float64(mEnd.TotalAlloc-mStart.TotalAlloc)/1024/1024, float64(mEnd.Sys-mStart.Sys)/1024/1024, mEnd.NumGC-mStart.NumGC)
	log.Infof("Execution Time: %v", elapsed)
	defer func() {
		err = os.Remove("./" + filename)
		if err != nil {
			log.Errorf("Failed to remove file %s: %v", filename, err)
		} else {
			log.Infof("Successfully removed file %s", filename)
		}
	}()
	return c.Download("./" + filename)
}
