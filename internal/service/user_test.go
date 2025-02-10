package service

import (
	"context"
	"github.com/dapoadeleke/balance-service/internal/http/dto"
	"github.com/dapoadeleke/balance-service/internal/model"
	"github.com/dapoadeleke/balance-service/internal/repository"
	"github.com/dapoadeleke/balance-service/internal/repository/mocks"
	"github.com/shopspring/decimal"
	"reflect"
	"testing"
)

func TestUserService_GetBalance(t *testing.T) {
	type fields struct {
		userRepository repository.User
	}
	type args struct {
		ctx context.Context
		id  uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    dto.UserBalanceResponse
		wantErr bool
	}{
		{
			name: "should return user balance",
			fields: fields{
				userRepository: func() repository.User {
					m := &mocks.User{}
					m.On("FindUserByID", context.TODO(), uint64(1)).Return(model.User{
						ID:      1,
						Name:    "James",
						Balance: decimal.NewFromInt(100),
					}, nil)
					return m
				}(),
			},
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			want: dto.UserBalanceResponse{
				UserID:  1,
				Balance: "100",
			},
			wantErr: false,
		},
		{
			name: "should return error when user not found",
			fields: fields{
				userRepository: func() repository.User {
					m := &mocks.User{}
					m.On("FindUserByID", context.TODO(), uint64(1)).Return(model.User{}, repository.ErrNoRecordFound)
					return m
				}(),
			},
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			want:    dto.UserBalanceResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserService{
				userRepository: tt.fields.userRepository,
			}
			got, err := u.GetBalance(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBalance() got = %v, want %v", got, tt.want)
			}
		})
	}
}
