package service

import (
	"database/sql"
	"fmt"

	"shiftmanager/internal/core/logger"
	"shiftmanager/internal/core/utils"
	"shiftmanager/internal/dto"
	"shiftmanager/internal/helper"
	"shiftmanager/internal/model"
	"shiftmanager/internal/repository"
)

type ShiftService interface {
	GetOne(input dto.ShiftPK) (dto.Shift, error)
	Save(input dto.SaveShift) error
	Generate(input dto.GenerateShift) error
}

type shiftService struct {
	shiftRepository          repository.ShiftRepository
	shiftPreferredRepository repository.ShiftPreferredRepository
	accountProfileRepository repository.AccountProfileRepository
}

func NewShiftService() ShiftService {
	return &shiftService{
		shiftRepository:          repository.NewShiftRepository(),
		shiftPreferredRepository: repository.NewShiftPreferredRepository(),
		accountProfileRepository: repository.NewAccountProfileRepository(),
	}
}

func (srv *shiftService) GetOne(input dto.ShiftPK) (dto.Shift, error) {
	if err := utils.ValidateYearMonth(input.Year, input.Month); err != nil {
		return dto.Shift{}, err
	}

	shift, err := srv.shiftRepository.GetOne(&model.Shift{
		Year:  input.Year,
		Month: input.Month,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			logger.Debug(err.Error())
			return dto.Shift{}, nil
		} else {
			logger.Error(err.Error())
			return dto.Shift{}, err
		}
	}

	var ret dto.Shift
	utils.MapFields(&ret, shift)
	return ret, nil
}

func (srv *shiftService) Save(input dto.SaveShift) error {
	if err := utils.ValidateYearMonth(input.Year, input.Month); err != nil {
		return err
	}
	holidays, err := utils.ParseDayCSV(input.StoreHoliday)
	if err != nil {
		return err
	}
	if err := utils.ValidateDaysInMonth(input.Year, input.Month, holidays); err != nil {
		return err
	}

	shift, err := srv.shiftRepository.GetOne(&model.Shift{
		Year:  input.Year,
		Month: input.Month,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			shift = model.Shift{
				Year:         input.Year,
				Month:        input.Month,
				StoreHoliday: input.StoreHoliday,
				Data:         input.Data,
			}
			err = srv.shiftRepository.Insert(&shift, nil)
		} else {
			logger.Error(err.Error())
			return err
		}
	} else {
		shift.StoreHoliday = input.StoreHoliday
		shift.Data = input.Data
		err = srv.shiftRepository.Update(&shift, nil)
	}

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (srv *shiftService) Generate(input dto.GenerateShift) error {
	if err := utils.ValidateYearMonth(input.Year, input.Month); err != nil {
		return err
	}
	holidays, err := utils.ParseDayCSV(input.StoreHoliday)
	if err != nil {
		return err
	}
	if err := utils.ValidateDaysInMonth(input.Year, input.Month, holidays); err != nil {
		return err
	}

	shiftGenerator := helper.NewShiftGenerator(input.Year, input.Month, holidays)
	err = shiftGenerator.InitRepositories()
	if err != nil {
		return fmt.Errorf("failed to initialize generator repositories: %v", err)
	}

	shiftCsv, err := shiftGenerator.GenerateCsvShift()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	err = srv.shiftRepository.Delete(&model.Shift{
		Year:  input.Year,
		Month: input.Month,
	}, nil)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	shift := model.Shift{
		Year:         input.Year,
		Month:        input.Month,
		StoreHoliday: input.StoreHoliday,
		Data:         &shiftCsv,
	}
	if err = srv.shiftRepository.Insert(&shift, nil); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}
