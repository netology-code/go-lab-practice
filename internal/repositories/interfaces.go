package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"github.com/yourname/finops-service/internal/models"
)

type UserRepository interface {
	GetBalance(ctx context.Context, userID int) (decimal.Decimal, error)
	UpdateBalance(ctx context.Context, userID int, amount decimal.Decimal, tx pgx.Tx) error
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *models.Transaction, dbTx pgx.Tx) (int, error)
	GetTransaction(ctx context.Context, id int) (*models.Transaction, error)
	UpdateTransaction(ctx context.Context, id int, newTx *models.Transaction) error
	DeleteTransaction(ctx context.Context, id int) error
}
