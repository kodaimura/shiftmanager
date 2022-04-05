package entity


type Shift struct {
	GId int `db: "GID json:"gid"`
	Year string `db: "YEAR json: "year"`
	Month string `db: "MONTH json: "month"`
	StoreHoliday string `db: STORE_HOLIDAY json: "storeholiday"`
	Shift string `db:"SHIFT" json:"shift"`
	CreateAt string `db:"CREATE_AT" json:"createat"`
	UpdateAt string `db:"UPDATE_AT" json:"updateat"`
}

