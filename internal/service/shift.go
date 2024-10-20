package service

import (
	"fmt"
	"strings"
	"time"
	"math"
	"math/rand"
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


func (srv *shiftService) getDailyPreferreds(year, month int, holidays []int) ([31][]int, error) {
	var ret [31][]int

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
		for _, date := range dates {
			ret[date - 1] = append(ret[date - 1], accountId)
		}
	}
	fmt.Println(ret)
	return ret, nil
}


func (srv *shiftService) getProfileMap() (map[int]model.AccountProfile, error) {
	ret := make(map[int]model.AccountProfile)

	profiles, err := srv.accountProfileRepository.Get(&model.AccountProfile{})
	if err != nil {
		return ret, err
	}

	for _, p := range profiles {
		ret[p.AccountId] = p
	}
	return ret, nil 
}


func (srv *shiftService) generateCsvShift(dailyPreferreds [31][]int, profileMap map[int]model.AccountProfile) string {
	var tmpShift [31][]int
	var shift [31][]int
	tmpScore := 0
	score := 0
	for _ = range 700000 {
		tmpShift = srv.makeTmpShift(dailyPreferreds)
		tmpScore = srv.evaluateShift(tmpShift, profileMap)
		if score < tmpScore {
			fmt.Println(tmpScore)
			fmt.Println(tmpShift)
			counts := make(map[int]int)
			for _, tmp := range shift {
				for _, x:= range tmp {
					counts[x]++
				}
			}
			fmt.Println(counts)
			score = tmpScore
			shift = tmpShift
		}
	}

	return srv.toCsvShift(shift, profileMap)
}


func (srv *shiftService) makeTmpShift(dailyPreferreds [31][]int) [31][]int {
	var tmp [31][]int
	for i := range 31 {
		tmp[i] = randomPick(dailyPreferreds[i])
	}
	return tmp
}


func randomPick(slice []int) []int {
	if len(slice) <= 2 {
		return slice
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})

	return slice[:2]
}


func (srv *shiftService) evaluateShift(shift [31][]int, profileMap map[int]model.AccountProfile) int {
	return srv.evaluateUniformity(shift) - srv.evaluateRolePenalty(shift, profileMap)
}


func (srv *shiftService) evaluateUniformity(shift [31][]int) int {
	counts := make(map[int]int)
	for _, tmp := range shift {
		for _, x:= range tmp {
			counts[x]++
		}
	}

	total := 0.0
	for _, num := range counts {
		total += float64(num)
	}

	mean := total / float64(len(counts))

	var varianceSum float64
	for _, num := range counts {
		varianceSum += math.Pow(float64(num)-mean, 2)
	}

	variance := varianceSum / float64(len(counts))
	stdDev := math.Sqrt(variance)

	if stdDev == 0 {
		return 100
	}
	return int(math.Min((1 / stdDev) * 100, 100))
}


func (srv *shiftService) evaluateRolePenalty(shift [31][]int, profileMap map[int]model.AccountProfile) int {
	penalty := 0
	for i := range 31 {
		if len(shift[i]) < 2 {
			continue
		}
		if profileMap[shift[i][0]].AccountRole == "1" && profileMap[shift[i][1]].AccountRole == "1" {
			// キッチン2人なら減点1
			penalty += 1
		}
		if profileMap[shift[i][0]].AccountRole == "2" && profileMap[shift[i][1]].AccountRole == "2" {
			// ホール2人なら減点2
			penalty += 2
		}
	}
	return penalty
}


func (srv *shiftService) toCsvShift(shift [31][]int, profileMap map[int]model.AccountProfile) string {
	ret := ""
	for i := range 31 {
		l := len(shift[i])
		if l == 0 {
			ret += ","
		}
		if l == 1 {
			ret += profileMap[shift[i][0]].DisplayName + ","
		}
		if l == 2 {
			ret += profileMap[shift[i][0]].DisplayName + " " + profileMap[shift[i][1]].DisplayName + ","
		}
	}
	return ret[:len(ret)-1]
}