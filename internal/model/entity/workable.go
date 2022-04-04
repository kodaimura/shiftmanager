package entity


type Workable struct {
	UId int `db:"UID" json:"uid"`
	Year string `db: "YEAR" json:"year"`
	Month string `db: "MONTH" json:"month"`
	WorkableDays string `db: WORKABLE_DAYS json:"workabledays"`
	Memo string `db: MEMO json:"memo"`
	CreateAt string `db:"CREATE_AT" json:"createat"`
	UpdateAt string `db:"UPDATE_AT" json:"updateat"`
}

//WorkableDays: 1,5,7,10,21,31
