package service

import (
	"context"
	"github.com/dapoadeleke/balance-service/internal/http/dto"
	"github.com/dapoadeleke/balance-service/internal/repository"
)

type User interface {
	GetBalance(ctx context.Context, id uint64) (dto.UserBalanceResponse, error)
}

type UserService struct {
	userRepository repository.User
}

func NewUserService(userRepository repository.User) *UserService {
	return &UserService{userRepository: userRepository}
}

func (u *UserService) GetBalance(ctx context.Context, id uint64) (dto.UserBalanceResponse, error) {
	user, err := u.userRepository.FindUserByID(ctx, id)
	if err != nil {
		return dto.UserBalanceResponse{}, err
	}

	return dto.UserBalanceResponse{
		UserID:  id,
		Balance: user.Balance.String(),
	}, nil
}
