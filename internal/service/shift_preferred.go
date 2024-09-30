package service

import (
	"database/sql"

	"shiftmanager/internal/core/logger"
	"shiftmanager/internal/core/utils"
	"shiftmanager/internal/dto"
	"shiftmanager/internal/model"
	"shiftmanager/internal/repository"
)

type ShiftPreferredService interface {
	GetOne(input dto.ShiftPreferredPK) (dto.ShiftPreferred, error)
	Save(input dto.SaveShiftPreferred) error
}

type shiftPreferredService struct {
	shiftPreferredRepository repository.ShiftPreferredRepository
}

func NewShiftPreferredService() ShiftPreferredService {
	return &shiftPreferredService{
		shiftPreferredRepository: repository.NewShiftPreferredRepository(),
	}
}

func (srv *shiftPreferredService) GetOne(input dto.ShiftPreferredPK) (dto.ShiftPreferred, error) {
	shiftPreferred, err := srv.shiftPreferredRepository.GetOne(&model.ShiftPreferred{
		AccountId: input.AccountId,
		Year: input.Year,
		Month: input.Month,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug(err.Error())
		} else {
			logger.Error(err.Error())
		}
		return dto.ShiftPreferred{}, err
	}

	var ret dto.ShiftPreferred
	utils.MapFields(&ret, shiftPreferred)
	return ret, nil
}


func (srv *shiftPreferredService) Save(input dto.SaveShiftPreferred) error {
	shiftPreferred, err := srv.shiftPreferredRepository.GetOne(&model.ShiftPreferred{ 
		AccountId: input.AccountId, 
		Year: input.Year, 
		Month: input.Month
	})

	if err != nil {
		if err == sql.ErrNoRows {
			shiftPreferred = &model.ShiftPreferred{
				AccountId: input.AccountId,
				Year: input.Year,
				Month: input.Month,
				Dates: input.Dates,
				Notes: input.Notes,
			}
			err = srv.shiftPreferredRepository.Insert(shiftPreferred, nil)
		} else {
			logger.Error(err.Error())
			return err
		}
	} else {
		shiftPreferred.Dates = input.Dates
		shiftPreferred.Notes = input.Notes
		err = srv.shiftPreferredRepository.Update(shiftPreferred, nil)
	}

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}