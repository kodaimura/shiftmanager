package dto


type Shift struct {
	Year int `json:"year"`
	Month int `json:"month"`
	StoreHoliday *string `json:"store_holiday"`
	Data *string `json:"shift_data"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ShiftPK struct {
	Year int `json:"year"`
	Month int `json:"month"`
}

type SaveShift struct {
	Year int `json:"year"`
	Month int `json:"month"`
	StoreHoliday *string `json:"store_holiday"`
	Data *string `json:"shift_data"`
}