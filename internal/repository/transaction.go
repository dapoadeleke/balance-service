package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dapoadeleke/balance-service/internal/db"
	"github.com/dapoadeleke/balance-service/internal/model"
	"time"
)

const (
	TransactionFieldTransactionID = "transaction_id"
	TransactionFieldUserID        = "user_id"
	TransactionFieldAmount        = "amount"
	TransactionFieldState         = "state"
	TransactionFieldSourceType    = "source_type"
	TransactionTableName          = "transactions"
)

type Transaction interface {
	SaveTransactionWithTx(ctx context.Context, tx db.Tx, transaction model.Transaction) (model.Transaction, error)
	FindTransactionByTransactionID(ctx context.Context, id string) (model.Transaction, error)
}

type TransactionRepository struct {
	db *db.Postgres
}

func NewTransactionRepository(db *db.Postgres) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) SaveTransactionWithTx(ctx context.Context, tx db.Tx, transaction model.Transaction) (model.Transaction, error) {
	fields := []string{
		TransactionFieldTransactionID,
		TransactionFieldUserID,
		TransactionFieldAmount,
		TransactionFieldState,
		TransactionFieldSourceType,
		CommonFieldCreatedAt,
		CommonFieldUpdatedAt,
	}

	now := time.Now()
	transaction.UpdatedAt = now

	var query string
	if transaction.ID == 0 {
		transaction.CreatedAt = now
		query = BuildSaveQuery(TransactionTableName, fields, false)
	} else {
		fields = append([]string{CommonFieldID}, fields...)
		query = BuildSaveQuery(TransactionTableName, fields, true)
	}

	_, err := tx.NamedExec(query, transaction)
	if err != nil {
		return model.Transaction{}, fmt.Errorf("failed to save transaction: %w", err)
	}

	return transaction, nil
}

func (r *TransactionRepository) FindTransactionByTransactionID(ctx context.Context, id string) (model.Transaction, error) {
	var t model.Transaction
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", TransactionTableName, TransactionFieldTransactionID)
	if err := r.db.Get(&t, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return t, ErrNoRecordFound
		}
		return t, fmt.Errorf("failed to get transaction by transation ID: %w", err)
	}
	return t, nil
}
