package request


type ShiftUri struct {
    Year  int `uri:"year" binding:"required"`
    Month int `uri:"month" binding:"required"`
}

type PostShift struct {
	StoreHoliday *string `json:"store_holiday"`
	Data *string `json:"shift_data"`
}

type PostShiftGenerate struct {
	StoreHoliday *string `json:"store_holiday"`
}