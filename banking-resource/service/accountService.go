package service

import (
	"banking/domain"
	"banking/dto"
	"banking/errs"
	"time"
)

const dbTSLayout = "2006-01-02 15:04:05"

type AccountService interface {
	NewAccount(dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError)
	MakeTransaction(dto.TransactionRequest) (*dto.TransactionResponse, *errs.AppError)
}

type DefaultAccountService struct {
	repo domain.AccountRespository
}

func (s DefaultAccountService) NewAccount(req dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError) {
	err := req.Validate()

	if err != nil {
		return nil, err
	}

	a := domain.Account{
		AccountId:   "",
		CustomerId:  req.CustomerId,
		OpeningDate: time.Now().Format("2006-01-02T15:04:05Z07:00"),
		AccountType: req.AccountType,
		Amount:      req.Amount,
		Status:      "1",
	}
	newAccount, err := s.repo.Save(a)

	if err != nil {
		return nil, err
	}
	response := newAccount.ToNewAccountResponseDto()

	return &response, nil
}

func (s DefaultAccountService) MakeTransaction(req dto.TransactionRequest) (*dto.TransactionResponse, *errs.AppError) {
	err := req.Validate()

	if err != nil {
		return nil, err
	}

	if req.IsTransactionTypeWithdrawal() {
		account, err := s.repo.FindBy(req.AccountId)
		if err != nil {
			return nil, err
		}
		if !account.CanWithdraw(req.Amount) {
			return nil, errs.NewValidationError("Insufficient balance in account")
		}
	}

	t := domain.Transaction{
		TransactionId:   "",
		AccountId:       req.AccountId,
		TransactionDate: time.Now().Format(dbTSLayout),
		TransactionType: req.TransactionType,
		Amount:          req.Amount,
	}
	newTransaction, err := s.repo.SaveTransaction(t)

	if err != nil {
		return nil, err
	}
	response := newTransaction.ToNewTransactionResponseDto()

	return &response, nil
}

func NewAccountService(repo domain.AccountRespository) DefaultAccountService {
	return DefaultAccountService{repo}
}
