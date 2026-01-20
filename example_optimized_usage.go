package main

import (
	"database/sql"
	"download-excel-project/export"
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

	// Demo export Excel dengan streaming
	fmt.Println("\n=== Demo Export Excel dengan Streaming ===")
	ExampleUsageStreamingExport()
}

// ExportToExcelWithStreaming - Contoh lengkap export Excel dengan streaming
func ExportToExcelWithStreaming(repo repository.CustomerRepo, totalData int32, filename string) error {
	fmt.Printf("Starting Excel export with streaming for %d records...\n", totalData)

	// Call our streaming export function
	resultFilename, err := export.ExportCustomerToExcelWithStreaming(repo, totalData, filename)
	if err != nil {
		return fmt.Errorf("error during Excel export: %v", err)
	}

	fmt.Printf("Excel export completed successfully: %s\n", resultFilename)
	return nil
}

// ExampleUsageStreamingExport - Contoh penggunaan export Excel dengan streaming
func ExampleUsageStreamingExport() {
	// Setup database connection
	db, err := sql.Open("postgres", "your-connection-string")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewCustomerRepo(db)

	// Export 100,000 records using streaming
	var totalData int32 = 100000
	filename := fmt.Sprintf("customer_updates_%s.xlsx", time.Now().Format("20060102_150405"))

	start := time.Now()
	err = ExportToExcelWithStreaming(repo, totalData, filename)
	if err != nil {
		log.Fatalf("Failed to export Excel: %v", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Total export time: %v\n", elapsed)
}

// DownloadCustomerExcelWithStreamingDemo - Fungsi demo untuk download Excel customer dengan streaming
func DownloadCustomerExcelWithStreamingDemo() {
	fmt.Println("=== Demo Download Excel Customer dengan Streaming ===")

	// Setup database connection (ganti dengan connection string yang sesuai)
	connectionString := "host=localhost port=5432 user=your_user password=your_password dbname=your_db sslmode=disable"

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		log.Println("Menggunakan mock data untuk demo...")
		// Untuk demo, kita bisa menggunakan mock repository
		mockRepo := &MockCustomerRepo{}
		demonstrateStreamingExport(mockRepo)
		return
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Printf("Cannot ping database: %v", err)
		log.Println("Menggunakan mock data untuk demo...")
		mockRepo := &MockCustomerRepo{}
		demonstrateStreamingExport(mockRepo)
		return
	}

	// Use real repository
	repo := repository.NewCustomerRepo(db)
	demonstrateStreamingExport(repo)
}

// demonstrateStreamingExport - Helper function untuk menjalankan demo export
func demonstrateStreamingExport(repo repository.CustomerRepo) {
	// Konfigurasi export
	var totalData int32 = 50000 // Export 50k records untuk demo
	filename := fmt.Sprintf("customer_updates_streaming_%s.xlsx", time.Now().Format("20060102_150405"))

	fmt.Printf("Memulai export %d records ke file: %s\n", totalData, filename)

	start := time.Now()
	err := ExportToExcelWithStreaming(repo, totalData, filename)
	if err != nil {
		log.Printf("Error during export: %v", err)
		return
	}

	elapsed := time.Since(start)
	fmt.Printf("Export selesai dalam waktu: %v\n", elapsed)
	fmt.Printf("File Excel berhasil dibuat: %s\n", filename)
}

// MockCustomerRepo - Mock repository untuk demo tanpa database
type MockCustomerRepo struct{}

func (m *MockCustomerRepo) GenerateUpdateData(totalData int32) ([][]interface{}, error) {
	// Mock implementation
	var data [][]interface{}
	for i := int32(0); i < totalData && i < 1000; i++ { // Limit to 1000 for demo
		row := []interface{}{
			fmt.Sprintf("2024-01-01 10:00:0%d", i%10),
			fmt.Sprintf("Customer %d", i+1),
			fmt.Sprintf("customer%d@example.com", i+1),
			fmt.Sprintf("REF%06d", i+1),
			"Profile Update",
		}
		data = append(data, row)
	}
	return data, nil
}

func (m *MockCustomerRepo) GenerateUpdateDataOptimized(batchSize int32, offset int32) ([][]interface{}, error) {
	return m.GenerateUpdateData(batchSize)
}

func (m *MockCustomerRepo) GenerateUpdateDataWithCursor(totalData int32) (<-chan []interface{}, <-chan error) {
	dataChan := make(chan []interface{}, 100)
	errChan := make(chan error, 1)

	go func() {
		defer close(dataChan)
		defer close(errChan)

		for i := int32(0); i < totalData; i++ {
			row := []interface{}{
				fmt.Sprintf("2024-01-01 10:%02d:%02d", i/60%60, i%60),
				fmt.Sprintf("Customer %d", i+1),
				fmt.Sprintf("customer%d@example.com", i+1),
				fmt.Sprintf("REF%08d", i+1),
				"Profile Update",
				fmt.Sprintf("+62812345%04d", i%10000),        // old_phone
				fmt.Sprintf("+62812346%04d", i%10000),        // new_phone
				fmt.Sprintf("old%d@example.com", i+1),        // old_email
				fmt.Sprintf("new%d@example.com", i+1),        // new_email
				"Engineer",                                   // old_occupation
				"Senior Engineer",                            // new_occupation
				"Software Developer",                         // old_job
				"Senior Software Developer",                  // new_job
				"PT ABC",                                     // old_company
				"PT XYZ",                                     // new_company
				"Jl. Old Company No. 123",                    // old_company_addr
				"Jl. New Company No. 456",                    // new_company_addr
				"Technology",                                 // old_company_desc
				"Technology Solutions",                       // new_company_desc
				"5-10 Million",                               // old_income
				"10-15 Million",                              // new_income
				fmt.Sprintf("+62813456%04d", i%10000),        // old_emergency_phone
				fmt.Sprintf("+62813457%04d", i%10000),        // new_emergency_phone
				fmt.Sprintf("Emergency Contact %d", i+1),     // old_emergency_contact
				fmt.Sprintf("New Emergency Contact %d", i+1), // new_emergency_contact
				"Jl. Old Address No. 123",                    // old_address
				"Jl. New Address No. 456",                    // new_address
				"Old Kecamatan",                              // old_kecamatan
				"New Kecamatan",                              // new_kecamatan
				"Old Kelurahan",                              // old_kelurahan
				"New Kelurahan",                              // new_kelurahan
				"001",                                        // old_rt
				"002",                                        // new_rt
				"003",                                        // old_rw
				"004",                                        // new_rw
				"Jakarta",                                    // old_city
				"Bekasi",                                     // new_city
				"DKI Jakarta",                                // old_province
				"Jawa Barat",                                 // new_province
				"12345",                                      // old_postal_code
				"54321",                                      // new_postal_code
			}
			dataChan <- row

			// Add small delay to simulate streaming
			if i%1000 == 0 {
				time.Sleep(10 * time.Millisecond)
				fmt.Printf("Generated %d records...\n", i+1)
			}
		}
	}()

	return dataChan, errChan
}

// main - Entry point untuk menjalankan demo
func demo() {
	fmt.Println("=== Demo Aplikasi Download Excel Customer dengan Streaming ===")
	fmt.Println()

	// Jalankan demo download Excel dengan streaming
	DownloadCustomerExcelWithStreamingDemo()

	fmt.Println()
	fmt.Println("=== Demo Selesai ===")
}
