package request


type PostAccountProfile struct {
	AccountRole string `json:"account_role"`
	DisplayName string `json:"display_name"`
}