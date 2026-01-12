package services

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/yourname/finops-service/internal/repositories"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetBalance(ctx context.Context, userID int) (decimal.Decimal, error) {
	// TODO: Реализовать получение баланса пользователя
	return decimal.Zero, nil
}
