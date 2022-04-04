package entity


type Group struct {
	GId int `db "GID" json: "gid"`
	GroupName string `db "GROUP_NAME" json: "groupname"`
	CreateAt string `db:"CREATE_AT" json:"createat"`
	UpdateAt string `db:"UPDATE_AT" json:"updateat"`
}