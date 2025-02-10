package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	ID            uint64                 `db:"id"`
	TransactionID string                 `db:"transaction_id"`
	UserID        uint64                 `db:"user_id"`
	Amount        decimal.Decimal        `db:"amount"`
	State         *TransactionState      `db:"state"`
	SourceType    *TransactionSourceType `db:"source_type"`
	CreatedAt     time.Time              `db:"created_at"`
	UpdatedAt     time.Time              `db:"updated_at"`
}

//go:generate go run github.com/dmarkham/enumer -type=TransactionState -json -sql -transform=upper -trimprefix=TransactionState -values
type TransactionState int

const (
	TransactionStateWin TransactionState = iota
	TransactionStateLose
)

//go:generate go run github.com/dmarkham/enumer -type=TransactionSourceType -json -sql -transform=upper -trimprefix=TransactionSourceType -values
type TransactionSourceType int

const (
	TransactionSourceTypeGame TransactionSourceType = iota
	TransactionSourceTypeServer
	TransactionSourceTypePayment
)
