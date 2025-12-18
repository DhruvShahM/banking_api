package repositories

import (
	"time"

	"banking-api/internal/models"

	"gorm.io/gorm"
)

type LoanPaymentRepository interface {
	Create(payment *models.LoanPayment) error
	GetByID(id int) (*models.LoanPayment, error)
	ListByLoan(loanID int) ([]models.LoanPayment, error)
	UpdateStatus(id int, status string, paidDate time.Time) error
}

type loanPaymentRepo struct {
	db *gorm.DB
}

func NewLoanPaymentRepo(db *gorm.DB) LoanPaymentRepository {
	return &loanPaymentRepo{db: db}
}

func (r *loanPaymentRepo) Create(payment *models.LoanPayment) error {
	return r.db.Create(payment).Error
}

func (r *loanPaymentRepo) GetByID(id int) (*models.LoanPayment, error) {
	var payment models.LoanPayment
	err := r.db.Preload("Loan").First(&payment, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

func (r *loanPaymentRepo) ListByLoan(loanID int) ([]models.LoanPayment, error) {
	var payments []models.LoanPayment
	err := r.db.Where("loan_id=?", loanID).Order("due_date ASC").Find(&payments).Error
	return payments, err
}

func (r *loanPaymentRepo) UpdateStatus(id int, status string, paidDate time.Time) error {
	return r.db.Model(&models.LoanPayment{}).Where("id=?", id).Updates(map[string]interface{}{
		"status":    status,
		"paid_date": paidDate,
	}).Error
}
