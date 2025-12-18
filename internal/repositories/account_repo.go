package repositories

import (
	"banking-api/internal/models"

	"gorm.io/gorm"
)

type accountRepo struct {
	db *gorm.DB
}

type AccountRepository interface {
	Create(account *models.Account) error
	GetByID(id int) (*models.Account, error)
	UpdateBalance(id int, amount float64) error
	ListByCustomer(customerID int) ([]models.Account, error)
}

func NewAccountRepo(db *gorm.DB) AccountRepository {
	return &accountRepo{db: db}
}

func (r *accountRepo) Create(account *models.Account) error {
	return r.db.Create(account).Error
}

func (r *accountRepo) GetByID(id int) (*models.Account, error) {
	var account models.Account
	err := r.db.Preload("Customer").First(&account, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *accountRepo) UpdateBalance(id int, amount float64) error {
	return r.db.Model(&models.Account{}).Where("id=?", id).Update("balance", gorm.Expr("balance+?", amount)).Error
}

func (r *accountRepo) ListByCustomer(customerID int) ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.Where("customer_id=?", customerID).Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, err
}
