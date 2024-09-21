package model


type AvailableWorkDays struct {
	AccountId int `db:"account_id" json:"account_id"`
	Year int `db:"year" json:"year"`
	Month int `db:"month" json:"month"`
	AvailableWorkDays *string `db:"available_work_days" json:"available_work_days"`
	Memo *string `db:"memo" json:"memo"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}