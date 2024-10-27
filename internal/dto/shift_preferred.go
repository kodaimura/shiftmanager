package dto


type ShiftPreferred struct {
	AccountId int `json:"account_id"`
	Year int `json:"year"`
	Month int `json:"month"`
	Dates *string `json:"dates"`
	Notes *string `json:"notes"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ShiftPreferredPK struct {
	AccountId int `json:"account_id"`
	Year int `json:"year"`
	Month int `json:"month"`
}

type SaveShiftPreferred struct {
	AccountId int `json:"account_id"`
	Year int `json:"year"`
	Month int `json:"month"`
	Dates *string `json:"dates"`
	Notes *string `json:"notes"`
}

type GetShiftPreferred struct {
	Year int `json:"year"`
	Month int `json:"month"`
}