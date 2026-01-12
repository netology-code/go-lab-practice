package services

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourname/finops-service/internal/models"
	"github.com/yourname/finops-service/internal/repositories"
)

type TransactionService struct {
	txRepo   repositories.TransactionRepository
	userRepo repositories.UserRepository
	pool     *pgxpool.Pool
}

func NewTransactionService(
	txRepo repositories.TransactionRepository,
	userRepo repositories.UserRepository,
	pool *pgxpool.Pool,
) *TransactionService {
	return &TransactionService{
		txRepo:   txRepo,
		userRepo: userRepo,
		pool:     pool,
	}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, tx *models.Transaction) (int, error) {
	// TODO: Реализовать создание транзакции с валидацией и управлением БД транзакцией
	return 0, nil
}

func (s *TransactionService) GetTransaction(ctx context.Context, id int) (*models.Transaction, error) {
	// TODO: Реализовать получение транзакции по ID
	return nil, nil
}

func (s *TransactionService) UpdateTransaction(ctx context.Context, id int, newTx *models.Transaction) error {
	// TODO: Реализовать обновление транзакции
	return nil
}

func (s *TransactionService) DeleteTransaction(ctx context.Context, id int) error {
	// TODO: Реализовать удаление транзакции
	return nil
}
