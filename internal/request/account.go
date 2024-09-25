package request


type Signup struct {
	Name string `json:"account_name"`
	Password string `json:"account_password"`
}

type Login struct {
	Name string `json:"account_name"`
	Password string `json:"account_password"`
}

type PutAccountPassword struct {
	OldPassword string `json:"old_account_password"`
	Password string `json:"account_password"`
}

type PutAccountName struct {
	Name string `json:"account_name"`
}