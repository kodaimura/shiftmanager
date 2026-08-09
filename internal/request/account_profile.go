package request

type PostAccountProfile struct {
	AccountRole string `json:"account_role" binding:"required,oneof=1 2 3"`
	DisplayName string `json:"display_name" binding:"required,max=1"`
}
