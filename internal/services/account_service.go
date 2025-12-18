package services

import (
	"banking-api/internal/models"
	"banking-api/internal/repositories"
	"errors"

	"gorm.io/gorm"
)

type AccountService interface {
	CreateAccount(req *models.CreateAccountRequest, customerID, branchID int) (*models.Account, error)
	Transfer(fromID, toID int, amount float64) error
	Deposit(accountID int, amount float64) error
	GetStatements(accountID int) ([]models.Transaction, error)
}

type accountService struct {
	db     *gorm.DB
	repo   repositories.AccountRepository
	txRepo repositories.TransctionRepository
}

func NewAccountService(db *gorm.DB, repo repositories.AccountRepository, txRepo repositories.TransctionRepository) AccountService {
	return &accountService{db: db, repo: repo, txRepo: txRepo}
}

func (s *accountService) CreateAccount(req *models.CreateAccountRequest, customerID, branchID int) (*models.Account, error) {
	account := &models.Account{
		CustomerID: customerID,
		BranchID:   branchID,
		Owner:      req.Owner,
		Currency:   req.Currency,
		Balance:    0,
	}

	if err := s.repo.Create(account); err != nil {
		return nil, err
	}
	return account, nil
}

func (s *accountService) Transfer(fromID, toID int, amount float64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		fromAcc, err := s.repo.GetByID(fromID)
		if err != nil {
			return errors.New("source account not found")
		}

		if fromAcc.Balance < amount {
			return errors.New("insufficient funds")
		}

		_, err = s.repo.GetByID(toID)
		if err != nil {
			return errors.New("destination account not found")
		}

		if err := s.repo.UpdateBalance(fromID, -amount); err != nil {
			return err
		}

		if err := s.repo.UpdateBalance(toID, amount); err != nil {
			return err
		}

		txn := &models.Transaction{
			FromAccountID: &fromID,
			ToAccountID:   &toID,
			Amount:        amount,
		}

		return s.txRepo.Create(txn)
	})
}

func (s *accountService) Deposit(accountID int, amount float64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		_, err := s.repo.GetByID(accountID)
		if err != nil {
			return errors.New("account not found")
		}

		if err := s.repo.UpdateBalance(accountID, amount); err != nil {
			return err
		}

		txn := &models.Transaction{
			FromAccountID: nil,
			ToAccountID:   &accountID,
			Amount:        amount,
		}

		return s.txRepo.Create(txn)
	})
}

func (s *accountService) GetStatements(accountID int) ([]models.Transaction, error) {
	var taxns []models.Transaction
	err := s.db.Where("from_account_id = ? OR to_account_id = ?", accountID, accountID).
		Order("created_at_DESC").Limit(100).Find(&taxns).Error
	return taxns, err
}
