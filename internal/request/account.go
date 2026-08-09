package request

type Signup struct {
	Name     string `json:"account_name" binding:"required,max=64"`
	Password string `json:"account_password" binding:"required,max=72"`
}

type Login struct {
	Name     string `json:"account_name" binding:"required,max=64"`
	Password string `json:"account_password" binding:"required,max=72"`
}

type PutAccountPassword struct {
	OldPassword string `json:"old_account_password" binding:"required,max=72"`
	Password    string `json:"account_password" binding:"required,max=72"`
}

type PutAccountName struct {
	Name string `json:"account_name" binding:"required,max=64"`
}
