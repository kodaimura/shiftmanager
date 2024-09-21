package model


type AccountProfile struct {
	AccountId int `db:"account_id" json:"account_id"`
	GroupId *int `db:"group_id" json:"group_id"`
	AccountRole string `db:"account_role" json:"account_role"`
	DisplayName *string `db:"display_name" json:"display_name"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}