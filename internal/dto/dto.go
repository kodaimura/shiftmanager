package dto


type WorkableExp1 struct {
	UId int `db:"UID" json:"uid"`
	Year string `db:"YEAR" json:"year"`
	Month string `db:"MONTH" json:"month"`
	WorkableDays string `db:"WORKABLE_DAYS" json:"workabledays"`
	Role string `db:"ROLE" json:"role"`
	AbbName string `db:"ABB_NAME" json:"abbname"`
}


type ProfileExp1 struct {
	UId int `db:"UID" json:"uid"`
	GId int `db:"GID" json:"gid"`
	Role string `db:"ROLE json:"role"`
	AbbName string `db:"ABB_NAME" json:"abbname"`
	GroupName string `db:"GROUP_NAME" json:"groupname"`
	RoleName string `db:"ROLE_NAME" json:"rolename"`
}