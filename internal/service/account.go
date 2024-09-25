package service

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"

	"shiftmanager/internal/core/jwt"
	"shiftmanager/internal/core/logger"
	"shiftmanager/internal/core/errs"
	"shiftmanager/internal/core/utils"
	"shiftmanager/internal/dto"
	"shiftmanager/internal/model"
	"shiftmanager/internal/repository"
)

type AccountService interface {
	GetOne(id int) (dto.Account, error)
	Delete(id int) error
	UpdateName(id int, name string) error
	UpdatePassword(id int, password string) error
	Login(input dto.Login) (dto.Account, error)
	Signup(input dto.Signup) (int, error)
	GenerateJwtPayload(id int) (jwt.Payload, error)
}

type accountService struct {
	accountRepository repository.AccountRepository
}

func NewAccountService() AccountService {
	return &accountService{
		accountRepository: repository.NewAccountRepository(),
	}
}

func (srv *accountService) GetOne(id int) (dto.Account, error) {
	account, err := srv.accountRepository.GetOne(&model.Account{Id: id})
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug(err.Error())
		} else {
			logger.Error(err.Error())
		}
		return dto.Account{}, err
	}

	var ret dto.Account
	utils.MapFields(&ret, account)
	return ret, nil
}

func (srv *accountService) UpdateName(id int, name string) error {
	if err := srv.checkUniqueName(id, name); err != nil {
		return err
	}

	account, err := srv.getAccountByID(id)
	if err != nil {
		return err
	}

	account.Name = name
	return srv.accountRepository.Update(account, nil)
}

func (srv *accountService) UpdatePassword(id int, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	account, err := srv.getAccountByID(id)
	if err != nil {
		return err
	}

	account.Password = string(hashed)
	return srv.accountRepository.Update(account, nil)
}

func (srv *accountService) Delete(id int) error {
	return srv.accountRepository.Delete(&model.Account{Id: id}, nil)
}

func (srv *accountService) Login(input dto.Login) (dto.Account, error) {
	account, err := srv.accountRepository.GetOne(&model.Account{Name: input.Name})
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug(err.Error())
		} else {
			logger.Error(err.Error())
		}
		return dto.Account{}, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(input.Password)); err != nil {
		logger.Error(err.Error())
		return dto.Account{}, err
	}

	var ret dto.Account
	utils.MapFields(&ret, account)
	return ret, nil
}

func (srv *accountService) Signup(input dto.Signup) (int, error) {
	if _, err := srv.accountRepository.GetOne(&model.Account{Name: input.Name}); err == nil {
		return 0, errs.NewUniqueConstraintError("account_name")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err.Error())
		return 0, err
	}

	account := &model.Account{
		Name:     input.Name,
		Password: string(hashed),
	}

	return srv.accountRepository.Insert(account, nil)
}

func (srv *accountService) GenerateJwtPayload(id int) (jwt.Payload, error) {
	account, err := srv.getAccountByID(id)
	if err != nil {
		return jwt.Payload{}, err
	}

	cc := jwt.CustomClaims{
		AccountId:   account.Id,
		AccountName: account.Name,
	}
	return jwt.NewPayload(cc), nil
}

func (srv *accountService) checkUniqueName(id int, name string) error {
	existingAccount, err := srv.accountRepository.GetOne(&model.Account{Name: name})
	if err == nil && existingAccount.Id != id {
		return errs.NewUniqueConstraintError("account_name")
	}
	return nil
}

func (srv *accountService) getAccountByID(id int) (*model.Account, error) {
	account, err := srv.accountRepository.GetOne(&model.Account{Id: id})
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return &account, nil
}
