package response


type GetAccount struct {
	Id int `json:"account_id"`
	Name string `json:"account_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type AccountPK struct {
	Id int `json:"account_id"`
}