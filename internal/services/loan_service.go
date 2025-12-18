package services

import (
	"banking-api/internal/models"
	"banking-api/internal/repositories"
	"math"
	"time"

	"gorm.io/gorm"
)

type LoanService interface {
	CreateLoan(req *models.CreateLoanRequest, customerID, brnachID int) (*models.Loan, error)
	ListLoans(customerID int) ([]models.Loan, error)
	UpdateStatus(loanID int, status string) error
}

type loanService struct {
	db     *gorm.DB
	repo   repositories.LoanRepository
	pmRepo repositories.LoanPaymentRepository
}

func NewLoanService(db *gorm.DB, repo repositories.LoanRepository, pmRepo repositories.LoanPaymentRepository) LoanService {
	return &loanService{db: db, repo: repo, pmRepo: pmRepo}
}

func calculateEMI(principal, monthlyRate float64, months int) float64 {
	if months == 0 {
		return 0
	}
	r := monthlyRate / 12

	if r == 0 {
		return principal / float64(months)
	}

	power := math.Pow(1+r, float64(months))
	emi := principal * r * power / (power - 1)
	return emi
}

func (s *loanService) CreateLoan(req *models.CreateLoanRequest, customerID, branchID int) (*models.Loan, error) {
	var loan *models.Loan
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// calculate total payable
		monthlyRate := req.InterestRate
		emi := calculateEMI(req.Amount, monthlyRate, req.TermMonths)
		totalPayable := emi * float64(req.TermMonths)

		loan = &models.Loan{
			CustomerID:   customerID,
			BranchID:     branchID,
			Amount:       req.Amount,
			InterestRate: req.InterestRate,
			TermMonths:   req.TermMonths,
			TotalPayable: totalPayable,
			Status:       "approved",
			StartDate:    time.Now(),
			EndDate:      time.Now().AddDate(0, req.TermMonths, 0),
		}

		if err := s.repo.Create(loan); err != nil {
			return err
		}

		// Generate payments (equal EMI)
		for i := 1; i <= req.TermMonths; i++ {
			dueDate := time.Now().AddDate(0, i, 0) // Monthly
			payment := &models.LoanPayment{
				LoanID:  loan.ID,
				Amount:  emi,
				DueDate: dueDate,
				Status:  "Pending",
			}

			if err := s.pmRepo.Create(payment); err != nil {
				return err // Rollback on failure
			}
		}

		// Reload with payments
		fullLoan, err := s.repo.GetByID(loan.ID)
		if err != nil {
			return err
		}

		loan = fullLoan
		return nil
	})

	if err != nil {
		return nil, err
	}

	return loan, nil
}

func (s *loanService) ListLoans(customerID int) ([]models.Loan, error) {
	return s.repo.ListByCustomer(customerID)
}

func (s *loanService) UpdateStatus(loanID int, status string) error {
	return s.repo.UpdateStatus(loanID, status)
}
