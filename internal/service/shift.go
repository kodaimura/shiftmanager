package service

import (
	"strings"
	"database/sql"

	"shiftmanager/internal/core/logger"
	"shiftmanager/internal/core/utils"
	"shiftmanager/internal/dto"
	"shiftmanager/internal/model"
	"shiftmanager/internal/repository"
)

type ShiftService interface {
	GetOne(input dto.ShiftPK) (dto.Shift, error)
	Save(input dto.SaveShift) error
	Generate(input dto.GenerateShift) error
}

type shiftService struct {
	shiftRepository repository.ShiftRepository
}

func NewShiftService() ShiftService {
	return &shiftService{
		shiftRepository: repository.NewShiftRepository(),
	}
}

func (srv *shiftService) GetOne(input dto.ShiftPK) (dto.Shift, error) {
	shift, err := srv.shiftRepository.GetOne(&model.Shift{
		Year: input.Year,
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
	shift, err := srv.shiftRepository.GetOne(&model.Shift{
		Year: input.Year, 
		Month: input.Month,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			shift = model.Shift{
				Year: input.Year,
				Month: input.Month,
				StoreHoliday: input.StoreHoliday,
				Data: input.Data,
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
	err := srv.shiftRepository.Delete(&model.Shift{
		Year: input.Year, 
		Month: input.Month,
	}, nil)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	shiftData, err := srv.generateProc(strings.Split(*input.StoreHoliday, ","))
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	shift := model.Shift{
		Year: input.Year,
		Month: input.Month,
		StoreHoliday: input.StoreHoliday,
		Data: &shiftData,
	}
	if err = srv.shiftRepository.Insert(&shift, nil); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}


func (srv *shiftService) generateProc(storeHoliday []string) (string, error) {
	return "", nil
}