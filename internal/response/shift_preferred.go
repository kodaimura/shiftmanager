package response


type GetShiftPreferred struct {
	AccountId int `json:"account_id"`
	Year int `json:"year"`
	Month int `json:"month"`
	Dates *string `json:"dates"`
	Notes *string `json:"notes"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}