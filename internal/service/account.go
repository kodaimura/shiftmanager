package service

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"

	"shiftmanager/internal/core/jwt"
	"shiftmanager/internal/core/logger"
	"shiftmanager/internal/core/errs"
	"shiftmanager/internal/model"
	"shiftmanager/internal/dto"
	"shiftmanager/internal/repository"
)


type AccountService interface {
	Signup(name, password string) error
	Login(name, password string) (dto.Account, error)
	GetProfile(id int) (dto.Account, error)
	GenerateJwtPayload(id int) (jwt.Payload, error)
	UpdateName(id int, name string) error
	UpdatePassword(id int, password string) error
	DeleteAccount(id int) error
}


type accountService struct {
	accountRepository repository.AccountRepository
}

func NewAccountService() AccountService {
	return &accountService{
		accountRepository: repository.NewAccountRepository(),
	}
}


func (srv *accountService) toAccountDTO(account model.Account) dto.Account {
	return dto.Account{
		Id:        account.Id,
		Name:      account.Name,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
}


func (srv *accountService) Signup(name, password string) error {
	_, err := srv.accountRepository.GetOne(&model.Account{Name: name})
	if err == nil {
		return errs.NewUniqueConstraintError("account_name")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	var account model.Account
	account.Name = name
	account.Password = string(hashed)

	err = srv.accountRepository.Insert(&account, nil);
	if err != nil {
		logger.Error(err.Error())
	}

	return err
}


func (srv *accountService) Login(name, password string) (dto.Account, error) {
	account, err := srv.accountRepository.GetOne(&model.Account{Name: name})
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug(err.Error())
		} else {
			logger.Error(err.Error())
		}
		return dto.Account{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		logger.Error(err.Error())
	}

	return srv.toAccountDTO(account), err
}


func (srv *accountService) GetProfile(id int) (dto.Account, error) {
	account, err := srv.accountRepository.GetOne(&model.Account{Id: id})
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug(err.Error())
		} else {
			logger.Error(err.Error())
		}
	}

	return srv.toAccountDTO(account), err
}


func (srv *accountService) GenerateJwtPayload(id int) (jwt.Payload, error) {
	account, err := srv.accountRepository.GetOne(&model.Account{Id: id})
	if err != nil {
		logger.Error(err.Error())
		return jwt.Payload{}, err
	}

	var cc jwt.CustomClaims
	cc.AccountId = account.Id
	cc.AccountName = account.Name
	return jwt.NewPayload(cc), nil
}


func (srv *accountService) UpdateName(id int, name string) error {
	u, err := srv.accountRepository.GetOne(&model.Account{Name: name})
	if err == nil && u.Id != id{
		return errs.NewUniqueConstraintError("account_name")
	}

	account, err := srv.accountRepository.GetOne(&model.Account{Id: id})
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	account.Name = name
	if err = srv.accountRepository.Update(&account, nil); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}


func (srv *accountService) UpdatePassword(id int, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	account, err := srv.accountRepository.GetOne(&model.Account{Id: id})
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	account.Password = string(hashed)
	if err = srv.accountRepository.Update(&account, nil); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}


func (srv *accountService) DeleteAccount(id int) error {
	if err := srv.accountRepository.Delete(&model.Account{Id: id}, nil); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}