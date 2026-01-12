package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) GetBalance(ctx context.Context, userID int) (decimal.Decimal, error) {
	// TODO: Реализовать получение баланса пользователя
	return decimal.Zero, nil
}

func (r *UserRepo) UpdateBalance(ctx context.Context, userID int, amount decimal.Decimal, tx pgx.Tx) error {
	// TODO: Реализовать обновление баланса пользователя
	return nil
}
