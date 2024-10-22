package helper

import (
	"math"
	"math/rand"
	"strings"
	"time"

	"shiftmanager/internal/core/utils"
	"shiftmanager/internal/model"
	"shiftmanager/internal/repository"
)

type ShiftGenerator struct {
	year          int
	month         int
	holidays      []int
	dailyPreferreds *[31][]int
	profileMap    map[int]model.AccountProfile
}

func NewShiftGenerator(year, month int, holidays []int) *ShiftGenerator {
	return &ShiftGenerator{
		year:   year,
		month:  month,
		holidays: holidays,
		profileMap: make(map[int]model.AccountProfile),
	}
}

func (gen *ShiftGenerator) InitRepositories() error {
	var err error
	gen.dailyPreferreds, err = gen.getDailyPreferreds()
	if err != nil {
		return err
	}

	for _, d := range gen.holidays {
		gen.dailyPreferreds[d-1] = []int{}
	}

	gen.profileMap, err = gen.getProfileMap()
	if err != nil {
		return err
	}
	return nil
}

func (gen *ShiftGenerator) getDailyPreferreds() (*[31][]int, error) {
	var ret [31][]int

	shiftPreferredRepo := repository.NewShiftPreferredRepository()
	preferreds, err := shiftPreferredRepo.Get(&model.ShiftPreferred{
		Year:  gen.year,
		Month: gen.month,
	})
	if err != nil {
		return &ret, err
	}

	for _, p := range preferreds {
		accountId := p.AccountId
		dates, _ := utils.AtoiSlice(strings.Split(*p.Dates, ","))
		for _, date := range dates {
			ret[date-1] = append(ret[date-1], accountId)
		}
	}
	return &ret, nil
}

func (gen *ShiftGenerator) getProfileMap() (map[int]model.AccountProfile, error) {
	ret := make(map[int]model.AccountProfile)

	profileRepo := repository.NewAccountProfileRepository()
	profiles, err := profileRepo.Get(&model.AccountProfile{})
	if err != nil {
		return ret, err
	}

	for _, p := range profiles {
		ret[p.AccountId] = p
	}
	return ret, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (gen *ShiftGenerator) GenerateCsvShift() (string, error) {
	shift := new([31][]int)
	tmpShift := new([31][]int)
	tmpScore := 0
	score := 0

	for i := 0; i < 700000; i++ {
		tmpShift = gen.makeTmpShift()
		tmpScore = gen.evaluateShift(tmpShift)
		if score < tmpScore {
			score = tmpScore
			*shift = *tmpShift
		}
	}

	return gen.toCsvShift(shift), nil
}

func (gen *ShiftGenerator) makeTmpShift() *[31][]int {
	tmp := new([31][]int)
	for i := 0; i < 31; i++ {
		tmp[i] = randomPick(gen.dailyPreferreds[i])
	}
	return tmp
}

func randomPick(slice []int) []int {
	if len(slice) <= 2 {
		return slice
	}
	first := rand.Intn(len(slice))
	second := rand.Intn(len(slice) - 1)
	if second >= first {
		second++
	}

	return []int{slice[first], slice[second]}
}

func (gen *ShiftGenerator) evaluateShift(shift *[31][]int) int {
	return gen.evaluateUniformity(shift) - gen.evaluateRolePenalty(shift)
}

func (gen *ShiftGenerator) evaluateUniformity(shift *[31][]int) int {
	counts := make(map[int]int)
	for _, tmp := range shift {
		for _, x := range tmp {
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
	return int(math.Min((1/stdDev)*100, 100))
}

func (gen *ShiftGenerator) evaluateRolePenalty(shift *[31][]int) int {
	penalty := 0
	for i := 0; i < 31; i++ {
		if len(shift[i]) < 2 {
			continue
		}
		if gen.profileMap[shift[i][0]].AccountRole == "1" && gen.profileMap[shift[i][1]].AccountRole == "1" {
			penalty += 1
		}
		if gen.profileMap[shift[i][0]].AccountRole == "2" && gen.profileMap[shift[i][1]].AccountRole == "2" {
			penalty += 2
		}
	}
	return penalty
}

func (gen *ShiftGenerator) toCsvShift(shift *[31][]int) string {
	ret := ""
	for i := 0; i < 31; i++ {
		l := len(shift[i])
		if l == 0 {
			ret += ","
		} else if l == 1 {
			ret += gen.profileMap[shift[i][0]].DisplayName + ","
		} else if l == 2 {
			ret += gen.profileMap[shift[i][0]].DisplayName + " " + gen.profileMap[shift[i][1]].DisplayName + ","
		}
	}
	return ret[:len(ret)-1]
}
