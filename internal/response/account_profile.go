package response


type GetAccountProfile struct {
	AccountId int `json:"account_id"`
	AccountRole string `json:"account_role"`
	DisplayName string `json:"display_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}