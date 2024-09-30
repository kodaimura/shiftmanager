package dto


type ShiftPreferred struct {
	AccountId string `json:"account_id"`
	Year int `json:"year"`
	Month int `json:"month"`
	Dates *string `json:"dates"`
	Notes *string `json:"notes"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ShiftPreferredPK struct {
	AccountId string `json:"account_id"`
	Year int `json:"year"`
	Month int `json:"month"`
}

type SaveShiftPreferred struct {
	AccountId string `json:"account_id"`
	Year int `json:"year"`
	Month int `json:"month"`
	Dates *string `json:"dates"`
	Notes *string `json:"notes"`
}