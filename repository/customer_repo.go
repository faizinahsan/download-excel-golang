package repository

import (
	"database/sql"
	"download-excel-project/domain"
	"encoding/json"
	"fmt"
)

type CustomerRepo interface {
	GenerateUpdateData(totalData int32) ([][]interface{}, error)
	GenerateUpdateDataOptimized(batchSize int32, offset int32) ([][]interface{}, error)
	GenerateUpdateDataWithCursor(totalData int32) (<-chan []interface{}, <-chan error)
}
type CustomerRepoImpl struct {
	db *sql.DB
}

func NewCustomerRepo(db *sql.DB) CustomerRepo {
	return &CustomerRepoImpl{db: db}
}

func (c CustomerRepoImpl) GenerateUpdateData(totalData int32) ([][]interface{}, error) {
	QueryGenerateUpdateData := `select
	tcup.id,
	tcup.customer_name,
	tcup.customer_email,
	tcup.created,
	tcup.reference_number,
	tcup.list_data_changes,
	tcup.old_data_mobile::jsonb->>'cellularNo' as old_phone,
	tcup.new_data_mobile::jsonb->>'cellularNo' as new_phone,
	tcup.old_data_mobile::jsonb->>'email' as old_email,
	tcup.new_data_mobile::jsonb->>'email' as new_email,
	tcup.old_data_mobile::jsonb->>'occupationDesc' as old_occupation,
	tcup.new_data_mobile::jsonb->>'occupationDesc' as new_occupation,
	tcup.old_data_mobile::jsonb->>'jobDesc' as old_job,
	tcup.new_data_mobile::jsonb->>'jobDesc' as new_job,
	tcup.old_data_mobile::jsonb->>'companyName' as old_companyName,
	tcup.new_data_mobile::jsonb->>'companyName' as new_companyName,
	tcup.old_data_mobile::jsonb->>'companyAddr' as old_companyAddr,
	tcup.new_data_mobile::jsonb->>'companyAddr' as new_companyAddr,
	tcup.old_data_mobile::jsonb->>'companyDesc' as old_companyDesc,
	tcup.new_data_mobile::jsonb->>'companyDesc' as new_companyDesc,
	tcup.old_data_mobile::jsonb->>'incomeDesc' as old_incomeDesc,
	tcup.new_data_mobile::jsonb->>'incomeDesc' as new_incomeDesc,
	tcup.old_data_mobile::jsonb->>'emergencyPhone' as old_emergencyPhone,
	tcup.new_data_mobile::jsonb->>'emergencyPhone' as new_emergencyPhone,
	tcup.old_data_mobile::jsonb->>'emergencyContact' as old_emergencyContact,
	tcup.new_data_mobile::jsonb->>'emergencyContact' as new_emergencyContact,
	tcup.old_data_mobile::jsonb->>'address' as old_street,
	tcup.new_data_mobile::jsonb->>'address' as new_street,
	tcup.old_data_mobile::jsonb->>'kecamatan' as old_kecamatan,
	tcup.new_data_mobile::jsonb->>'kecamatan' as new_kecamatan,
	tcup.old_data_mobile::jsonb->>'kelurahan' as old_kelurahan,
	tcup.new_data_mobile::jsonb->>'kelurahan' as new_kelurahan,
	tcup.old_data_mobile::jsonb->>'rt' as old_rt,
	tcup.new_data_mobile::jsonb->>'rt' as new_rt,
	tcup.old_data_mobile::jsonb->>'rw' as old_rw,
	tcup.new_data_mobile::jsonb->>'rw' as new_rw,
	tcup.old_data_mobile::jsonb->>'cityName' as old_cityName,
	tcup.new_data_mobile::jsonb->>'cityName' as new_cityName,
	tcup.old_data_mobile::jsonb->>'provinceName' as old_provinceName,
	tcup.new_data_mobile::jsonb->>'provinceName' as new_provinceName,
	tcup.old_data_mobile::jsonb->>'postalCode' as old_postalCode,
	tcup.new_data_mobile::jsonb->>'postalCode' as new_postalCode
from
	customer_update_profile tcup ORDER BY tcup.id DESC LIMIT $1`
	rows, err := c.db.Query(QueryGenerateUpdateData, totalData)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers [][]interface{}
	for rows.Next() {
		var customer domain.UpdateProfile
		err := rows.Scan(&customer.Id, &customer.Name, &customer.Email, &customer.Created, &customer.ReferenceNumber, &customer.ListDataChanges,
			&customer.OldPhone, &customer.NewPhone,
			&customer.OldEmail, &customer.NewEmail,
			&customer.OldOccupation, &customer.NewOccupation,
			&customer.OldJob, &customer.NewJob,
			&customer.OldCompanyName, &customer.NewCompanyName,
			&customer.OldCompanyAddr, &customer.NewCompanyAddr,
			&customer.OldCompanyDesc, &customer.NewCompanyDesc,
			&customer.OldIncomeDesc, &customer.NewIncomeDesc,
			&customer.OldEmergencyPhone, &customer.NewEmergencyPhone,
			&customer.OldEmergencyContact, &customer.NewEmergencyContact,
			&customer.OldStreet, &customer.NewStreet,
			&customer.OldKecamatan, &customer.NewKecamatan,
			&customer.OldKelurahan, &customer.NewKelurahan,
			&customer.OldRT, &customer.NewRT,
			&customer.OldRW, &customer.NewRW,
			&customer.OldCityName, &customer.NewCityName,
			&customer.OldProvinceName, &customer.NewProvinceName,
			&customer.OldPostalCode, &customer.NewPostalCode)
		if err != nil {
			return nil, err
		}
		row := []interface{}{
			customer.Created,
			customer.Name,
			customer.Email,
			customer.ReferenceNumber,
			customer.ListDataChanges,
			customer.OldPhone,
			customer.NewPhone,
			customer.OldEmail,
			customer.NewEmail,
			customer.OldOccupation,
			customer.NewOccupation,
			customer.OldJob,
			customer.NewJob,
			customer.OldCompanyName,
			customer.NewCompanyName,
			customer.OldCompanyAddr,
			customer.NewCompanyAddr,
			customer.OldCompanyDesc,
			customer.NewCompanyDesc,
			customer.OldIncomeDesc,
			customer.NewIncomeDesc,
			customer.OldEmergencyPhone,
			customer.NewEmergencyPhone,
			customer.OldEmergencyContact,
			customer.NewEmergencyContact,
			customer.OldStreet,
			customer.NewStreet,
			customer.OldKecamatan,
			customer.NewKecamatan,
			customer.OldKelurahan,
			customer.NewKelurahan,
			customer.OldRT,
			customer.NewRT,
			customer.OldRW,
			customer.NewRW,
			customer.OldCityName,
			customer.NewCityName,
			customer.OldProvinceName,
			customer.NewProvinceName,
			customer.OldPostalCode,
			customer.NewPostalCode,
		}
		customers = append(customers, row)
	}
	return customers, nil
}

// GenerateUpdateDataOptimized - Optimasi dengan batch processing
func (c CustomerRepoImpl) GenerateUpdateDataOptimized(batchSize int32, offset int32) ([][]interface{}, error) {
	// Query yang lebih optimal - ambil raw JSON dan parse di Go
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

		// Parse JSON di Go side (lebih efisien)
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

// GenerateUpdateDataWithCursor - Streaming data dengan cursor untuk memory efficiency
func (c CustomerRepoImpl) GenerateUpdateDataWithCursor(totalData int32) (<-chan []interface{}, <-chan error) {
	dataChan := make(chan []interface{}, 1000) // Buffer untuk 1000 records
	errChan := make(chan error, 1)

	go func() {
		defer close(dataChan)
		defer close(errChan)

		var batchSize int32 = 10000 // Process 10k records per batch
		var offset int32 = 0

		for offset < totalData {
			currentBatch := batchSize
			if offset+batchSize > totalData {
				currentBatch = totalData - offset
			}

			fmt.Printf("Processing batch: offset=%d, size=%d\n", offset, currentBatch)

			batch, err := c.GenerateUpdateDataOptimized(currentBatch, offset)
			if err != nil {
				errChan <- err
				return
			}

			for _, row := range batch {
				dataChan <- row
			}

			offset += batchSize
		}
	}()

	return dataChan, errChan
}
