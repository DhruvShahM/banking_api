package services

import (
	"errors"
	"time"

	"banking-api/internal/models"
	"banking-api/internal/repositories"

	"gorm.io/gorm"
)

type LoanPaymentService interface {
	MakePayment(paymentID int, loanID int) error
	ListPayments(loanID int) ([]models.LoanPayment, error)
}

type loanPaymentService struct {
	db       *gorm.DB
	repo     repositories.LoanPaymentRepository
	loanRepo repositories.LoanRepository
}

func NewLoanPaymentService(db *gorm.DB, repo repositories.LoanPaymentRepository, loanRepo repositories.LoanRepository) LoanPaymentService {
	return &loanPaymentService{db: db, repo: repo, loanRepo: loanRepo}
}

func (s *loanPaymentService) MakePayment(paymentID, loanID int) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		payment, err := s.repo.GetByID(paymentID)
		if err != nil {
			return errors.New("payment not found")
		}

		if payment.LoanID != loanID {
			return errors.New("payment does not belong to this loan")
		}

		if payment.Status != "paid" {
			return errors.New("payment already made")
		}

		// Update Payment
		if err := s.repo.UpdateStatus(paymentID, "paid", time.Now()); err != nil {
			return err
		}

		// UPDATE LOAN IF ALL PAID
		loan, err := s.loanRepo.GetByID(loanID)
		if err != nil {
			return errors.New("loan not found")
		}
		payments, _ := s.repo.ListByLoan(loanID)
		paidCount := 0
		for _, p := range payments {
			if p.Status == "paid" {
				paidCount++
			}
		}

		if paidCount == loan.TermMonths {
			s.loanRepo.UpdateStatus(loanID, "repaid")
		}

		return nil
	})
}

func (s *loanPaymentService) ListPayments(loanID int) ([]models.LoanPayment, error) {
	return s.repo.ListByLoan(loanID)
}
