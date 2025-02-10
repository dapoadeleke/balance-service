package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type User struct {
	ID        uint64          `json:"id" db:"id"`
	Name      string          `json:"name" db:"name"`
	Balance   decimal.Decimal `json:"balance" db:"balance"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}
