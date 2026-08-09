package request

type ShiftUri struct {
	Year  int `uri:"year" binding:"required,gte=2000,lte=2100"`
	Month int `uri:"month" binding:"required,gte=1,lte=12"`
}

type PostShift struct {
	StoreHoliday *string `json:"store_holiday" binding:"omitempty,max=100"`
	Data         *string `json:"shift_data" binding:"omitempty,max=1000"`
}

type PostShiftGenerate struct {
	StoreHoliday *string `json:"store_holiday" binding:"omitempty,max=100"`
}
