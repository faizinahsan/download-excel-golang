package main

import (
	"database/sql"
	"download-excel-project/repository"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func optimized() {
	// Contoh penggunaan method optimized

	// Setup database connection
	db, err := sql.Open("postgres", "your-connection-string")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewCustomerRepo(db)

	// Method 1: Batch Processing - Sangat direkomendasikan untuk 1 juta data
	fmt.Println("=== Method 1: Batch Processing ===")
	start := time.Now()

	var allData [][]interface{}
	var totalData int32 = 1000000
	var batchSize int32 = 50000 // Process 50k records per batch
	var offset int32 = 0

	for offset < totalData {
		currentBatch := batchSize
		if offset+batchSize > totalData {
			currentBatch = totalData - offset
		}

		fmt.Printf("Processing batch: offset=%d, size=%d\n", offset, currentBatch)

		batch, err := repo.GenerateUpdateDataOptimized(currentBatch, offset)
		if err != nil {
			log.Fatal(err)
		}

		allData = append(allData, batch...)
		offset += batchSize
	}

	elapsed1 := time.Since(start)
	fmt.Printf("Batch processing selesai: %d records dalam %v\n", len(allData), elapsed1)

	// Method 2: Streaming dengan Channel - Paling memory efficient
	fmt.Println("\n=== Method 2: Streaming dengan Channel ===")
	start = time.Now()

	dataChan, errChan := repo.GenerateUpdateDataWithCursor(totalData)
	var count int

	go func() {
		for row := range dataChan {
			count++
			// Process each row here
			_ = row // Placeholder untuk processing

			if count%10000 == 0 {
				fmt.Printf("Processed %d records so far...\n", count)
			}
		}
	}()

	// Check for errors
	select {
	case err := <-errChan:
		if err != nil {
			log.Fatal(err)
		}
	case <-time.After(10 * time.Minute): // Timeout after 10 minutes
		fmt.Println("Processing completed or timed out")
	}

	elapsed2 := time.Since(start)
	fmt.Printf("Streaming selesai: %d records dalam %v\n", count, elapsed2)

	// Method 3: Original (untuk perbandingan) - TIDAK direkomendasikan untuk 1 juta data
	fmt.Println("\n=== Method 3: Original Method (untuk perbandingan) ===")
	start = time.Now()

	// HATI-HATI: Method ini bisa menggunakan banyak memory untuk 1 juta data
	// Sebaiknya test dengan data yang lebih sedikit dulu
	originalData, err := repo.GenerateUpdateData(10000) // Test dengan 10k saja
	if err != nil {
		log.Fatal(err)
	}

	elapsed3 := time.Since(start)
	fmt.Printf("Original method selesai: %d records dalam %v\n", len(originalData), elapsed3)

	// Perbandingan performa
	fmt.Println("\n=== Perbandingan Performa ===")
	fmt.Printf("Batch Processing: %v\n", elapsed1)
	fmt.Printf("Streaming: %v\n", elapsed2)
	fmt.Printf("Original (10k records): %v\n", elapsed3)
}

// Contoh untuk Export Excel dengan Streaming
func ExportToExcelWithStreaming(repo repository.CustomerRepo, totalData int32) error {
	dataChan, errChan := repo.GenerateUpdateDataWithCursor(totalData)

	// Setup Excel file (contoh dengan excelize)
	// f := excelize.NewFile()
	// defer f.Close()

	var rowNum = 2 // Start from row 2 (row 1 untuk header)

	for {
		select {
		case _, ok := <-dataChan:
			if !ok {
				// Channel closed, processing completed
				fmt.Printf("Excel export selesai: total %d rows\n", rowNum-2)
				return nil
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
				return fmt.Errorf("error during streaming: %v", err)
			}
		}
	}
}
