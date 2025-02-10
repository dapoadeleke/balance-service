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
	UserFieldName    = "name"
	UserFieldBalance = "balance"
	UserTableName    = "users"
)

type User interface {
	SaveUserWithTx(ctx context.Context, tx db.Tx, user model.User) (model.User, error)
	FindUserByID(ctx context.Context, id uint64) (model.User, error)
}

type UserRepository struct {
	db *db.Postgres
}

func NewUserRepository(db *db.Postgres) *UserRepository {
	return &UserRepository{db: db}
}

func (p *UserRepository) SaveUserWithTx(ctx context.Context, tx db.Tx, user model.User) (model.User, error) {
	fields := []string{
		UserFieldName,
		UserFieldBalance,
		CommonFieldCreatedAt,
		CommonFieldUpdatedAt,
	}

	now := time.Now()
	user.UpdatedAt = now

	var query string
	if user.ID == 0 {
		user.CreatedAt = now
		query = BuildSaveQuery(TransactionTableName, fields, false)
	} else {
		fields = append([]string{CommonFieldID}, fields...)
		query = BuildSaveQuery(UserTableName, fields, true)
	}

	_, err := tx.NamedExec(query, user)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) FindUserByID(ctx context.Context, id uint64) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", UserTableName, CommonFieldID)
	if err := r.db.Get(&user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrNoRecordFound
		}
		return user, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
