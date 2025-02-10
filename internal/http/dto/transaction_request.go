package dto

import (
	"errors"
	"github.com/dapoadeleke/balance-service/internal/model"
	"github.com/shopspring/decimal"
	"strconv"
)

type TransactionRequest struct {
	UserID        string `json:"userId"`
	Amount        string `json:"amount"`
	TransactionID string `json:"transactionId"`
	State         string
	SourceType    string
}

func (t *TransactionRequest) Validate() error {
	if t.UserID == "" {
		return errors.New("invalid user id")
	}
	if t.State == "" {
		return errors.New("state is required")
	}
	if t.Amount == "" {
		return errors.New("amount is required")
	}
	if amount, err := decimal.NewFromString(t.Amount); err != nil {
		return errors.New("invalid amount")
	} else if amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be greater than 0")
	}
	if t.TransactionID == "" {
		return errors.New("transaction id is required")
	}

	return nil
}

func (t *TransactionRequest) ToTransaction() (model.Transaction, error) {
	userID, err := strconv.ParseUint(t.UserID, 10, 64)
	if err != nil {
		return model.Transaction{}, err
	}

	amount, err := decimal.NewFromString(t.Amount)
	if err != nil {
		return model.Transaction{}, err
	}

	state, err := model.TransactionStateString(t.State)
	if err != nil {
		return model.Transaction{}, errors.New("invalid state")
	}

	sourceType, err := model.TransactionSourceTypeString(t.SourceType)
	if err != nil {
		return model.Transaction{}, errors.New("invalid source type")
	}

	return model.Transaction{
		UserID:        userID,
		Amount:        amount,
		TransactionID: t.TransactionID,
		State:         &state,
		SourceType:    &sourceType,
	}, nil
}
