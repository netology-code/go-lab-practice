package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourname/finops-service/internal/models"
)

type TransactionRepo struct {
	pool *pgxpool.Pool
}

func NewTransactionRepo(pool *pgxpool.Pool) *TransactionRepo {
	return &TransactionRepo{pool: pool}
}

func (r *TransactionRepo) CreateTransaction(ctx context.Context, tx *models.Transaction, dbTx pgx.Tx) (int, error) {
	// TODO: Реализовать создание транзакции в БД
	return 0, nil
}

func (r *TransactionRepo) GetTransaction(ctx context.Context, id int) (*models.Transaction, error) {
	// TODO: Реализовать получение транзакции по ID
	return nil, nil
}

func (r *TransactionRepo) UpdateTransaction(ctx context.Context, id int, newTx *models.Transaction) error {
	// TODO: Реализовать обновление транзакции
	return nil
}

func (r *TransactionRepo) DeleteTransaction(ctx context.Context, id int) error {
	// TODO: Реализовать удаление транзакции
	return nil
}
