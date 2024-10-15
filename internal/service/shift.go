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
	shiftPreferredRepository repository.ShiftPreferredRepository
	accountProfileRepository repository.AccountProfileRepository
}

func NewShiftService() ShiftService {
	return &shiftService{
		shiftRepository: repository.NewShiftRepository(),
		shiftPreferredRepository: repository.NewShiftPreferredRepository(),
		accountProfileRepository: repository.NewAccountProfileRepository(),
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

	shiftData, err := srv.generateProc(input)
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


func (srv *shiftService) generateProc(input dto.GenerateShift) (string, error) {
	holidays, _ := utils.AtoiSlice(strings.Split(*input.StoreHoliday, ","))
	dailyPreferreds, err := srv.getDailyPreferreds(input.Year, input.Month, holidays)
	if err != nil {
		return "", err
	}

	profileMap, err := srv.getProfileMap()
	if err != nil {
		return "", err
	}
	
	return srv.generateCsvShift(dailyPreferreds, profileMap), nil
}


func (srv *shiftService) getDailyPreferreds(year, month int, holidays []int) ([32][]int, error) {
	var ret [32][]int

	preferreds, err := srv.shiftPreferredRepository.Get(&model.ShiftPreferred{
		Year: year,
		Month: month,
	})
	if err != nil {
		return ret, err
	}
	
	for _, p := range preferreds {
		accountId := p.AccountId
		dates, _ := utils.AtoiSlice(strings.Split(*p.Dates, ","))
		for date := range dates {
			ret[date] = append(ret[date], accountId)
		}
	}
	return ret, nil
}


func (srv *shiftService) getProfileMap() (map[int]model.AccountProfile, error) {
	var ret map[int]model.AccountProfile

	profiles, err := srv.accountProfileRepository.Get(&model.AccountProfile{})
	if err != nil {
		return ret, err
	}

	for _, p := range profiles {
		ret[p.AccountId] = p
	}
	return ret, nil 
}


func (srv *shiftService) generateCsvShift(dailyPreferreds [32][]int, profileMap map[int]model.AccountProfile) string {
	var tmpShift [32][]int
	var shift [32][]int
	var tmpScore int
	var score int
	for _ = range 10000 {
		tmpShift = srv.makeTmpShift(dailyPreferreds)
		tmpScore = srv.evaluateShift(tmpShift, profileMap)
		if score < tmpScore {
			score = tmpScore
			shift = tmpShift
		}
	}

	return srv.toCsv(shift, profileMap)
}


func (srv *shiftService) makeTmpShift(dailyPreferreds [32][]int) [32][]int {	
	return [32][]int{}
}


func (srv *shiftService) evaluateShift(shift [32][]int, profileMap map[int]model.AccountProfile) int {	
	return 0
}


func (srv *shiftService) toCsv(shift [32][]int, profileMap map[int]model.AccountProfile) string {	
	return ""
}