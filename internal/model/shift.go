package model


type Shift struct {
	Year int `db:"year" json:"year"`
	Month int `db:"month" json:"month"`
	StoreHoliday *string `db:"store_holiday" json:"store_holiday"`
	Data *string `db:"shift_data" json:"shift_data"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}