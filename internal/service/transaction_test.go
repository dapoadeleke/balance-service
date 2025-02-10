package service

import (
	"context"
	"errors"
	"github.com/dapoadeleke/balance-service/internal/db"
	"github.com/dapoadeleke/balance-service/internal/db/mocks"
	"github.com/dapoadeleke/balance-service/internal/model"
	"github.com/dapoadeleke/balance-service/internal/repository"
	mockRepo "github.com/dapoadeleke/balance-service/internal/repository/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestTransactionService_PostTransaction(t1 *testing.T) {
	type fields struct {
		db                    db.DB
		userRepository        repository.User
		transactionRepository repository.Transaction
	}
	type args struct {
		ctx         context.Context
		transaction model.Transaction
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "should post transaction successfully",
			fields: fields{
				db: func() db.DB {
					d := &mocks.DB{}
					tx := &mocks.Tx{}
					tx.On("Commit").Return(nil)
					tx.On("Rollback").Return(nil)
					d.On("MustBegin").Return(tx, nil)
					return d
				}(),
				userRepository: func() repository.User {
					u := &mockRepo.User{}
					u.On("FindUserByID", context.TODO(), uint64(1)).Return(model.User{
						ID:      1,
						Name:    "John",
						Balance: decimal.NewFromInt(100),
					}, nil)
					u.On("SaveUserWithTx", context.TODO(), mock.AnythingOfType("*mocks.Tx"), mock.Anything).Return(model.User{}, nil)
					return u
				}(),
				transactionRepository: func() repository.Transaction {
					t := &mockRepo.Transaction{}
					t.On("FindTransactionByTransactionID", context.TODO(), "TXN_001").Return(model.Transaction{}, repository.ErrNoRecordFound)
					t.On("SaveTransactionWithTx", context.TODO(), mock.AnythingOfType("*mocks.Tx"), mock.Anything).Return(model.Transaction{}, nil)
					return t
				}(),
			},
			args: args{
				ctx: context.TODO(),
				transaction: model.Transaction{
					TransactionID: "TXN_001",
					UserID:        1,
					Amount:        decimal.NewFromInt(100),
					State:         ptr(model.TransactionStateWin),
					SourceType:    ptr(model.TransactionSourceTypeGame),
				},
			},
			wantErr: false,
		}, {
			name: "should fail when user not found",
			fields: fields{
				db: &mocks.DB{},
				userRepository: func() repository.User {
					u := &mockRepo.User{}
					u.On("FindUserByID", context.TODO(), uint64(1)).Return(model.User{}, repository.ErrNoRecordFound)
					return u
				}(),
			},
			args: args{
				ctx: context.TODO(),
				transaction: model.Transaction{
					TransactionID: "TXN_001",
					UserID:        1,
					Amount:        decimal.NewFromInt(100),
					State:         ptr(model.TransactionStateWin),
					SourceType:    ptr(model.TransactionSourceTypeGame),
				},
			},
			wantErr: true,
		},
		{
			name: "should fail when transaction already processed",
			fields: fields{
				db: &mocks.DB{},
				userRepository: func() repository.User {
					u := &mockRepo.User{}
					u.On("FindUserByID", context.TODO(), uint64(1)).Return(model.User{
						ID:      1,
						Name:    "John",
						Balance: decimal.NewFromInt(100),
					}, nil)
					u.On("SaveUserWithTx", context.TODO(), mock.AnythingOfType("*mocks.Tx"), mock.Anything).Return(model.User{}, nil)
					return u
				}(),
				transactionRepository: func() repository.Transaction {
					t := &mockRepo.Transaction{}
					t.On("FindTransactionByTransactionID", context.TODO(), "TXN_001").Return(model.Transaction{
						TransactionID: "TXN_001",
					}, nil)
					return t
				}(),
			},
			args: args{
				ctx: context.TODO(),
				transaction: model.Transaction{
					TransactionID: "TXN_001",
					UserID:        1,
					Amount:        decimal.NewFromInt(100),
					State:         ptr(model.TransactionStateWin),
					SourceType:    ptr(model.TransactionSourceTypeGame),
				},
			},
			wantErr: true,
		},
		{
			name: "should fail when balance is not sufficient",
			fields: fields{
				db: &mocks.DB{},
				userRepository: func() repository.User {
					u := &mockRepo.User{}
					u.On("FindUserByID", context.TODO(), uint64(1)).Return(model.User{
						ID:      1,
						Name:    "John",
						Balance: decimal.NewFromInt(100),
					}, nil)
					return u
				}(),
				transactionRepository: func() repository.Transaction {
					t := &mockRepo.Transaction{}
					t.On("FindTransactionByTransactionID", context.TODO(), "TXN_001").Return(model.Transaction{}, repository.ErrNoRecordFound)
					return t
				}(),
			},
			args: args{
				ctx: context.TODO(),
				transaction: model.Transaction{
					TransactionID: "TXN_001",
					UserID:        1,
					Amount:        decimal.NewFromInt(150),
					State:         ptr(model.TransactionStateLose),
					SourceType:    ptr(model.TransactionSourceTypeGame),
				},
			},
			wantErr: true,
		},
		{
			name: "should fail when database error occurs while saving user",
			fields: fields{
				db: func() db.DB {
					d := &mocks.DB{}
					tx := &mocks.Tx{}
					tx.On("Commit").Return(nil)
					tx.On("Rollback").Return(nil)
					d.On("MustBegin").Return(tx, nil)
					return d
				}(),
				userRepository: func() repository.User {
					u := &mockRepo.User{}
					u.On("FindUserByID", context.TODO(), uint64(1)).Return(model.User{
						ID:      1,
						Name:    "John",
						Balance: decimal.NewFromInt(100),
					}, nil)
					u.On("SaveUserWithTx", context.TODO(), mock.AnythingOfType("*mocks.Tx"), mock.Anything).Return(model.User{}, errors.New("failed to save user"))
					return u
				}(),
				transactionRepository: func() repository.Transaction {
					t := &mockRepo.Transaction{}
					t.On("FindTransactionByTransactionID", context.TODO(), "TXN_001").Return(model.Transaction{}, repository.ErrNoRecordFound)
					t.On("SaveTransactionWithTx", context.TODO(), mock.AnythingOfType("*mocks.Tx"), mock.Anything).Return(model.Transaction{}, nil)
					return t
				}(),
			},
			args: args{
				ctx: context.TODO(),
				transaction: model.Transaction{
					TransactionID: "TXN_001",
					UserID:        1,
					Amount:        decimal.NewFromInt(100),
					State:         ptr(model.TransactionStateWin),
					SourceType:    ptr(model.TransactionSourceTypeGame),
				},
			},
			wantErr: true,
		},
		{
			name: "should fail when database error occurs while saving transaction",
			fields: fields{
				db: func() db.DB {
					d := &mocks.DB{}
					tx := &mocks.Tx{}
					tx.On("Commit").Return(nil)
					tx.On("Rollback").Return(nil)
					d.On("MustBegin").Return(tx, nil)
					return d
				}(),
				userRepository: func() repository.User {
					u := &mockRepo.User{}
					u.On("FindUserByID", context.TODO(), uint64(1)).Return(model.User{
						ID:      1,
						Name:    "John",
						Balance: decimal.NewFromInt(100),
					}, nil)
					u.On("SaveUserWithTx", context.TODO(), mock.AnythingOfType("*mocks.Tx"), mock.Anything).Return(model.User{}, nil)
					return u
				}(),
				transactionRepository: func() repository.Transaction {
					t := &mockRepo.Transaction{}
					t.On("FindTransactionByTransactionID", context.TODO(), "TXN_001").Return(model.Transaction{}, repository.ErrNoRecordFound)
					t.On("SaveTransactionWithTx", context.TODO(), mock.AnythingOfType("*mocks.Tx"), mock.Anything).Return(model.Transaction{}, errors.New("failed to save transaction"))
					return t
				}(),
			},
			args: args{
				ctx: context.TODO(),
				transaction: model.Transaction{
					TransactionID: "TXN_001",
					UserID:        1,
					Amount:        decimal.NewFromInt(100),
					State:         ptr(model.TransactionStateWin),
					SourceType:    ptr(model.TransactionSourceTypeGame),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TransactionService{
				db:                    tt.fields.db,
				userRepository:        tt.fields.userRepository,
				transactionRepository: tt.fields.transactionRepository,
			}
			if err := t.PostTransaction(tt.args.ctx, tt.args.transaction); (err != nil) != tt.wantErr {
				t1.Errorf("PostTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
