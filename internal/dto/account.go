package dto


type Account struct {
	Id int `json:"account_id"`
	Name string `json:"account_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Signup struct {
	Name string `json:"account_name"`
	Password string `json:"account_password"`
}

type Login struct {
	Name string `json:"account_name"`
	Password string `json:"account_password"`
}