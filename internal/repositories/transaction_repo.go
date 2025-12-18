package repositories

import (
	"banking-api/internal/models"

	"gorm.io/gorm"
)

type TransctionRepository interface {
	Create(txn *models.Transaction) error
}

type transactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) TransctionRepository {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) Create(txn *models.Transaction) error {
	return r.db.Create(txn).Error
}
