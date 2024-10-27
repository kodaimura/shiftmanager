package model


type Account struct {
	Id int `db:"account_id" json:"account_id"`
	Name string `db:"account_name" json:"account_name"`
	Password string `db:"account_password" json:"account_password"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}