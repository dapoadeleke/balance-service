package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/dapoadeleke/balance-service/internal/db"
	"github.com/dapoadeleke/balance-service/internal/model"
	"github.com/dapoadeleke/balance-service/internal/repository"
)

type Transaction interface {
	PostTransaction(ctx context.Context, transaction model.Transaction) error
}

type TransactionService struct {
	db                    db.DB
	userRepository        repository.User
	transactionRepository repository.Transaction
}

func NewTransactionService(
	db db.DB,
	userRepository repository.User,
	transactionRepository repository.Transaction,
) *TransactionService {
	return &TransactionService{
		db:                    db,
		userRepository:        userRepository,
		transactionRepository: transactionRepository,
	}
}

func (t *TransactionService) PostTransaction(ctx context.Context, transaction model.Transaction) error {
	user, err := t.userRepository.FindUserByID(ctx, transaction.UserID)
	if err != nil {
		return fmt.Errorf("failed to find user by id: %w", err)
	}

	_, err = t.transactionRepository.FindTransactionByTransactionID(ctx, transaction.TransactionID)
	if err == nil {
		return fmt.Errorf("transaction ID %s has already been processed", transaction.TransactionID)
	} else if !errors.Is(err, repository.ErrNoRecordFound) {
		return fmt.Errorf("failed to find transaction by transaction id: %w", err)
	}

	switch *transaction.State {
	case model.TransactionStateWin:
		user.Balance = user.Balance.Add(transaction.Amount)
	case model.TransactionStateLose:
		if user.Balance.LessThan(transaction.Amount) {
			return fmt.Errorf("insufficient balance")
		}
		user.Balance = user.Balance.Sub(transaction.Amount)
	default:
		return fmt.Errorf("invalid transaction state: %s", transaction.State)
	}

	tx := t.db.MustBegin()
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	_, err = t.transactionRepository.SaveTransactionWithTx(ctx, tx, transaction)
	if err != nil {
		return fmt.Errorf("failed to save transaction with transaction: %w", err)
	}

	_, err = t.userRepository.SaveUserWithTx(ctx, tx, user)
	if err != nil {
		return fmt.Errorf("failed to save user with transaction: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
