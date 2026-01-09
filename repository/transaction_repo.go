package repository

import (
	"database/sql"
	"download-excel-project/domain"
)

type TransactionRepo interface {
	GetTransactions() ([]domain.Transaction, error)
	GetTransactionToArray() ([][]interface{}, error)
}

type TransactionRepoImpl struct {
	db *sql.DB
}

func NewTransactionRepo(db *sql.DB) TransactionRepo {
	return &TransactionRepoImpl{db: db}
}

func (t TransactionRepoImpl) GetTransactions() ([]domain.Transaction, error) {
	rows, err := t.db.Query("SELECT id, transaction_amount, transaction_date, feature_name, reference_number FROM t_transaction order by transaction_date DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var transaction domain.Transaction
		err := rows.Scan(&transaction.ID, &transaction.Amount,
			&transaction.Date, &transaction.TransactionName, &transaction.ReferenceNumber)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (t TransactionRepoImpl) GetTransactionToArray() ([][]interface{}, error) {
	rows, err := t.db.Query("SELECT id, transaction_amount, transaction_date, feature_name, reference_number FROM transactions order by id, transaction_date DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions [][]interface{}
	for rows.Next() {
		var transaction domain.Transaction
		err := rows.Scan(&transaction.ID, &transaction.Amount,
			&transaction.Date, &transaction.TransactionName, &transaction.ReferenceNumber)
		if err != nil {
			return nil, err
		}
		row := []interface{}{
			transaction.ID,
			transaction.Amount,
			transaction.Date.Format("2006-01-02"),
			transaction.TransactionName,
			transaction.ReferenceNumber,
		}
		transactions = append(transactions, row)
	}
	return transactions, nil
}
