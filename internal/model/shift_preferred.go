package model


type ShiftPreferred struct {
	AccountId int `db:"account_id" json:"account_id"`
	Year int `db:"year" json:"year"`
	Month int `db:"month" json:"month"`
	Dates *string `db:"dates" json:"dates"`
	Notes *string `db:"notes" json:"notes"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}