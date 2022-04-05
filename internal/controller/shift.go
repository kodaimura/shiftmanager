package controller

import (
    "fmt"
    "log"
    "strings"
    "math/rand"
    "github.com/gin-gonic/gin"
    
    "shiftmanager/internal/pkg/jwtauth"
    "shiftmanager/pkg/utils"
    "shiftmanager/internal/dto"
    "shiftmanager/internal/model/repository"
    "shiftmanager/internal/model/entity"
    "shiftmanager/internal/constants"
)


type ShiftController interface {
    ShiftPage(c *gin.Context)
    ShiftGenerate(c *gin.Context)
    Shift(c *gin.Context)
    StoreHoliday(c *gin.Context)
}


type shiftController struct {
    ur repository.UserRepository
    wr repository.WorkableRepository
    sr repository.ShiftRepository
}


func NewShiftController() ShiftController {
    ur := repository.NewUserRepository()
    wr := repository.NewWorkableRepository()
    sr := repository.NewShiftRepository()
    return shiftController{ur, wr, sr}
}


//GET /shift/:year/:month
func (sc shiftController)ShiftPage(c *gin.Context) {
    gid, err := jwtauth.ExtractGId(c)

    if err != nil {
        c.Redirect(303, "/login")
        return
    }

    year := c.Param("year")
    month := c.Param("month")

    shift, err := sc.sr.SelectByGIdYearMonth(gid, year, month)

    if err != nil {
        c.Redirect(303, "/shift/" + year + "/" + month + "/storeholiday")
        return
    }

    workabledays, _ := sc.wr.GetWorkableExp1ByGId(gid, year, month)
    base0 := makeBase0(&workabledays)

    c.HTML(200, "shift.html", gin.H{
        "appname": constants.AppName,
        "year": year,
        "month": month,
        "shift": shift.Shift,
        "storeholiday": shift.StoreHoliday,
        "candidate": shapeShiftToRender(&base0),
    })
}


//GET /shift/:year/:month/storeholiday
func (sc shiftController)StoreHoliday(c *gin.Context) {
    year := c.Param("year")
    month := c.Param("month")

    c.HTML(200, "storeholiday.html", gin.H{
        "appname": constants.AppName,
        "year": year,
        "month": month,
    })
}


//POST /shift/:year/:month
func (sc shiftController)Shift(c *gin.Context) {
    gid, err := jwtauth.ExtractGId(c)

    if err != nil {
        c.Redirect(303, "/login")
        return
    }
    year := c.Param("year")
    month := c.Param("month")

    s := &entity.Shift{}
    s.Year = year
    s.Month = month
    s.StoreHoliday = c.PostForm("storeholiday")
    s.Shift = c.PostForm("shift")
    s.GId = gid

    err = sc.sr.Upsert(*s)

    if err != nil {
        fmt.Println(err)
    }

    c.Redirect(303, "/shift/" + year + "/" + month)
}


//POST /shift/:year/:month/auto-generate
func (sc shiftController)ShiftGenerate(c *gin.Context) {
    gid, err := jwtauth.ExtractGId(c)

    if err != nil {
        c.Redirect(303, "/login")
        return
    }

    year := c.Param("year")
    month := c.Param("month")
    storeholiday := c.PostForm("storeholiday")
    workabledays, _ := sc.wr.GetWorkableExp1ByGId(gid, year, month)

    roleMap := makeRoleMap(&workabledays)
    base0 := makeBase0(&workabledays)
    shift := makeShift(&base0, &roleMap, storeholiday)

    c.HTML(200, "shift.html", gin.H{
        "appname": constants.AppName,
        "shift": shapeShiftToRender(&shift),
        "storeholiday": storeholiday,
        "candidate": shapeShiftToRender(&base0),
        "year": year,
        "month": month,
    })
}


func shapeShiftToRender(shift *[31][]string) string {
    csv := ""
    for _, s := range *shift {
        csv += ","
        for i, name := range s {
            if i > 0 {
                csv += " "
            }
            csv += name
        }
    }
    return csv[1:]
}


func makeShift(base0 *[31][]string, roleMap *map[string]string, csvHoliday string) [31][]string {
    holiday, _ := utils.AtoiSlice(strings.Split(csvHoliday, ","))

    reflectHolidayToBase0(holiday, base0)

    var base [31][][]string = makeBase(base0)
    var lens [31]int = LensForEachBaseElem(&base)

    return pickShift(&base, &lens, roleMap)
}


func makeRoleMap(data *[]dto.WorkableExp1) map[string]string {
    roleMap := map[string]string{}
    for _, d := range *data {
        roleMap[d.AbbName] = d.Role
     }

     return roleMap
}


func makeBase0(data *[]dto.WorkableExp1) [31][]string {
    var base0 [31][]string

    for _, d := range *data {
        days, err := utils.AtoiSlice(strings.Split(d.WorkableDays, ","))
        if err != nil {
            log.Panic(err)
            return [31][]string{}
        }

        for _, day := range days {
            base0[day - 1] = append(base0[day - 1], d.AbbName) 
        }        
     }
     return base0
}


func reflectHolidayToBase0(holiday []int, base0 *[31][]string) {
    for _, day := range holiday {
        base0[day - 1] = []string{}
    }
}


func makeBase(base0 *[31][]string) [31][][]string{
    var base [31][][]string

    for i, sl := range base0 {
        base[i] = utils.Combinations(sl, 2)
    }
    return base
}


func LensForEachBaseElem(base *[31][][]string) [31]int{
    var lens [31]int

    for i, sl := range *base {
        lens[i] = len(sl)
    }
    return lens
}


func pickShift(base *[31][][]string, lens *[31]int, roleMap *map[string]string) [31][]string{
    var shift, rs [31][]string
    var max = -99999.9
    var p = 0.0

    for i := 0; i < 1000000; i++ {
        rs = makeRandomShift(base, lens)
        p = evaluateShift(&rs, roleMap)

        if max < p {
            shift = rs
            max = p
        }
    }
    return shift
}


func makeRandomShift(base *[31][][]string, lens *[31]int) [31][]string{
    var rs [31][]string
    var r int

    for i, sl := range *base {
        r = rand.Intn(lens[i])
        rs[i] = sl[r]
    }

    return rs
}


func evaluateShift(rs *[31][]string, roleMap *map[string]string) float64 {
    return rolePoint(rs, roleMap) - variancePoint(rs, len(*roleMap))
}


func rolePoint(rs *[31][]string, roleMap *map[string]string) float64 {
    point := 0.0
    r1, r2, r3 := 0, 0, 0

    for _, names := range rs {
        r1, r2, r3 = 0, 0, 0
        for _, n := range names {
            switch (*roleMap)[n] {
                case "1":
                    r1 += 1
                case "2":
                    r2 += 1
                case "3":
                    r3 += 1
            }
        }

        if r1 == 1 && r2 == 1 {
            point += 1.5
        } else if r3 >= 1 {
            point += 1.2
        } else if r1 == 1 {
            point += 1.0
        }
    }
    return point
}


func variancePoint(rs *[31][]string, numOfMember int) float64 {
    nom := float64(numOfMember)
    sl := flatten(rs)
    m := countEachElems(&sl)
    sum, sumSqu := 0.0, 0.0

    for _, count := range m {
        sum += float64(count)
        sumSqu += float64(count * count)
    }

    return (sumSqu / nom) - ((sum / nom) * (sum / nom))
}


func flatten(sl *[31][]string) []string {
    var fsl []string
    
    for _, s := range *sl {
        fsl = append(fsl, s...)
    }
    return fsl
}


func countEachElems(sl *[]string) map[string]int {
    m := map[string]int{}

    for _, e := range *sl {
        m[e] += 1
    }
    return m
}