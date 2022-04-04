package entity


type General struct {
	Class string `db "CLASS" json: "class"`
	Key1 string `db "KEY1" json: "key1"`
	Key2 string `db "KEY2" json: "key2"`
	Value1 string `db:"VALUE1" json:"value1"`
	Value2 string `db:"VALUE2" json:"value2"`
	Remarks string `db:"REMARKS" json:"remarks"`
	CreateAt string `db:"CREATE_AT" json:"createat"`
	UpdateAt string `db:"UPDATE_AT" json:"updateat"`
}