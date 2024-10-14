package request


type PostShift struct {
	StoreHoliday *string `json:"store_holiday"`
	Data *string `json:"data"`
}

type ShiftUri struct {
    Year  int `uri:"year" binding:"required"`
    Month int `uri:"month" binding:"required"`
}