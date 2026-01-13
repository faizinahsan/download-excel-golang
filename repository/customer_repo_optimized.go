package repository

import (
	"database/sql"
	"encoding/json"
)

type CustomerRepoOptimized struct {
	db *sql.DB
}

func NewCustomerRepoOptimized(db *sql.DB) *CustomerRepoOptimized {
	return &CustomerRepoOptimized{db: db}
}

// Optimasi 1: Menggunakan Batch Processing dengan Cursor
func (c *CustomerRepoOptimized) GenerateUpdateDataWithBatch(batchSize int32, offset int32) ([][]interface{}, error) {
	// Query yang lebih optimal dengan OFFSET dan LIMIT yang lebih kecil
	query := `
	SELECT 
		tcup.id,
		tcup.customer_name,
		tcup.customer_email,
		tcup.created,
		tcup.reference_number,
		tcup.list_data_changes,
		tcup.old_data_mobile,
		tcup.new_data_mobile
	FROM customer_update_profile tcup 
	ORDER BY tcup.id DESC 
	OFFSET $1 LIMIT $2`

	rows, err := c.db.Query(query, offset, batchSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers [][]interface{}
	for rows.Next() {
		var customer struct {
			Id              int32
			Name            string
			Email           string
			Created         string
			ReferenceNumber string
			ListDataChanges string
			OldDataMobile   string
			NewDataMobile   string
		}

		err := rows.Scan(
			&customer.Id,
			&customer.Name,
			&customer.Email,
			&customer.Created,
			&customer.ReferenceNumber,
			&customer.ListDataChanges,
			&customer.OldDataMobile,
			&customer.NewDataMobile,
		)
		if err != nil {
			return nil, err
		}

		// Parse JSON di Go side (lebih efisien dari PostgreSQL)
		oldData := make(map[string]interface{})
		newData := make(map[string]interface{})

		json.Unmarshal([]byte(customer.OldDataMobile), &oldData)
		json.Unmarshal([]byte(customer.NewDataMobile), &newData)

		// Helper function to safely get string from map
		getString := func(data map[string]interface{}, key string) string {
			if val, ok := data[key]; ok {
				if str, ok := val.(string); ok {
					return str
				}
			}
			return ""
		}

		row := []interface{}{
			customer.Created,
			customer.Name,
			customer.Email,
			customer.ReferenceNumber,
			customer.ListDataChanges,
			getString(oldData, "cellularNo"),       // old_phone
			getString(newData, "cellularNo"),       // new_phone
			getString(oldData, "email"),            // old_email
			getString(newData, "email"),            // new_email
			getString(oldData, "occupationDesc"),   // old_occupation
			getString(newData, "occupationDesc"),   // new_occupation
			getString(oldData, "jobDesc"),          // old_job
			getString(newData, "jobDesc"),          // new_job
			getString(oldData, "companyName"),      // old_companyName
			getString(newData, "companyName"),      // new_companyName
			getString(oldData, "companyAddr"),      // old_companyAddr
			getString(newData, "companyAddr"),      // new_companyAddr
			getString(oldData, "companyDesc"),      // old_companyDesc
			getString(newData, "companyDesc"),      // new_companyDesc
			getString(oldData, "incomeDesc"),       // old_incomeDesc
			getString(newData, "incomeDesc"),       // new_incomeDesc
			getString(oldData, "emergencyPhone"),   // old_emergencyPhone
			getString(newData, "emergencyPhone"),   // new_emergencyPhone
			getString(oldData, "emergencyContact"), // old_emergencyContact
			getString(newData, "emergencyContact"), // new_emergencyContact
			getString(oldData, "address"),          // old_street
			getString(newData, "address"),          // new_street
			getString(oldData, "kecamatan"),        // old_kecamatan
			getString(newData, "kecamatan"),        // new_kecamatan
			getString(oldData, "kelurahan"),        // old_kelurahan
			getString(newData, "kelurahan"),        // new_kelurahan
			getString(oldData, "rt"),               // old_rt
			getString(newData, "rt"),               // new_rt
			getString(oldData, "rw"),               // old_rw
			getString(newData, "rw"),               // new_rw
			getString(oldData, "cityName"),         // old_cityName
			getString(newData, "cityName"),         // new_cityName
			getString(oldData, "provinceName"),     // old_provinceName
			getString(newData, "provinceName"),     // new_provinceName
			getString(oldData, "postalCode"),       // old_postalCode
			getString(newData, "postalCode"),       // new_postalCode
		}
		customers = append(customers, row)
	}
	return customers, nil
}

// Optimasi 2: Menggunakan Streaming dengan Channel untuk Memory Efficiency
func (c *CustomerRepoOptimized) StreamUpdateData(totalData int32, batchSize int32, dataChan chan<- []interface{}) error {
	defer close(dataChan)

	var offset int32 = 0
	for offset < totalData {
		currentBatch := batchSize
		if offset+batchSize > totalData {
			currentBatch = totalData - offset
		}

		batch, err := c.GenerateUpdateDataWithBatch(currentBatch, offset)
		if err != nil {
			return err
		}

		for _, row := range batch {
			dataChan <- row
		}

		offset += batchSize
	}
	return nil
}

// Optimasi 3: Query dengan Index Hint (jika menggunakan PostgreSQL)
func (c *CustomerRepoOptimized) GenerateUpdateDataWithIndexOptimization(totalData int32) ([][]interface{}, error) {
	// Menggunakan WHERE condition untuk membatasi range ID
	query := `
	WITH max_id AS (
		SELECT MAX(id) as max_id FROM customer_update_profile
	)
	SELECT 
		tcup.id,
		tcup.customer_name,
		tcup.customer_email,
		tcup.created,
		tcup.reference_number,
		tcup.list_data_changes,
		tcup.old_data_mobile,
		tcup.new_data_mobile
	FROM customer_update_profile tcup, max_id
	WHERE tcup.id > (max_id.max_id - $1)
	ORDER BY tcup.id DESC`

	rows, err := c.db.Query(query, totalData)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers [][]interface{}
	for rows.Next() {
		var customer struct {
			Id              int32
			Name            string
			Email           string
			Created         string
			ReferenceNumber string
			ListDataChanges string
			OldDataMobile   string
			NewDataMobile   string
		}

		err := rows.Scan(
			&customer.Id,
			&customer.Name,
			&customer.Email,
			&customer.Created,
			&customer.ReferenceNumber,
			&customer.ListDataChanges,
			&customer.OldDataMobile,
			&customer.NewDataMobile,
		)
		if err != nil {
			return nil, err
		}

		// Parse JSON di Go side
		oldData := make(map[string]interface{})
		newData := make(map[string]interface{})

		json.Unmarshal([]byte(customer.OldDataMobile), &oldData)
		json.Unmarshal([]byte(customer.NewDataMobile), &newData)

		getString := func(data map[string]interface{}, key string) string {
			if val, ok := data[key]; ok {
				if str, ok := val.(string); ok {
					return str
				}
			}
			return ""
		}

		row := []interface{}{
			customer.Created,
			customer.Name,
			customer.Email,
			customer.ReferenceNumber,
			customer.ListDataChanges,
			getString(oldData, "cellularNo"),
			getString(newData, "cellularNo"),
			getString(oldData, "email"),
			getString(newData, "email"),
			getString(oldData, "occupationDesc"),
			getString(newData, "occupationDesc"),
			getString(oldData, "jobDesc"),
			getString(newData, "jobDesc"),
			getString(oldData, "companyName"),
			getString(newData, "companyName"),
			getString(oldData, "companyAddr"),
			getString(newData, "companyAddr"),
			getString(oldData, "companyDesc"),
			getString(newData, "companyDesc"),
			getString(oldData, "incomeDesc"),
			getString(newData, "incomeDesc"),
			getString(oldData, "emergencyPhone"),
			getString(newData, "emergencyPhone"),
			getString(oldData, "emergencyContact"),
			getString(newData, "emergencyContact"),
			getString(oldData, "address"),
			getString(newData, "address"),
			getString(oldData, "kecamatan"),
			getString(newData, "kecamatan"),
			getString(oldData, "kelurahan"),
			getString(newData, "kelurahan"),
			getString(oldData, "rt"),
			getString(newData, "rt"),
			getString(oldData, "rw"),
			getString(newData, "rw"),
			getString(oldData, "cityName"),
			getString(newData, "cityName"),
			getString(oldData, "provinceName"),
			getString(newData, "provinceName"),
			getString(oldData, "postalCode"),
			getString(newData, "postalCode"),
		}
		customers = append(customers, row)
	}
	return customers, nil
}
