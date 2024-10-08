package service

import (
	"database/sql"

	"shiftmanager/internal/core/logger"
	"shiftmanager/internal/core/utils"
	"shiftmanager/internal/dto"
	"shiftmanager/internal/model"
	"shiftmanager/internal/repository"
)

type AccountProfileService interface {
	GetOne(accountId int) (dto.AccountProfile, error)
	Save(accountId int, input dto.SaveAccountProfile) error
}

type accountProfileService struct {
	accountProfileRepository repository.AccountProfileRepository
}

func NewAccountProfileService() AccountProfileService {
	return &accountProfileService{
		accountProfileRepository: repository.NewAccountProfileRepository(),
	}
}

func (srv *accountProfileService) GetOne(accountId int) (dto.AccountProfile, error) {
	accountProfile, err := srv.accountProfileRepository.GetOne(&model.AccountProfile{AccountId: accountId})
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug(err.Error())
			return dto.AccountProfile{}, nil
		} else {
			logger.Error(err.Error())
			return dto.AccountProfile{}, err
		}
	}

	var ret dto.AccountProfile
	utils.MapFields(&ret, accountProfile)
	return ret, nil
}


func (srv *accountProfileService) Save(accountId int, input dto.SaveAccountProfile) error {
	accountProfile := &model.AccountProfile{
		AccountId:   accountId,
		AccountRole: input.AccountRole,
		DisplayName: input.DisplayName,
	}

	_, err := srv.accountProfileRepository.GetOne(&model.AccountProfile{AccountId: accountId})
	if err == nil {
		err = srv.accountProfileRepository.Update(accountProfile, nil)
	} else {
		err = srv.accountProfileRepository.Insert(accountProfile, nil)
	}

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}