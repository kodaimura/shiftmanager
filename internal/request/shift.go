package request


type PostShift struct {
	StoreHoliday *string `json:"store_holiday"`
	Data *string `json:"data"`
}

type ShiftUri struct {
    Year  int `form:"year" binding:"required"`
    Month int `form:"month" binding:"required"`
}