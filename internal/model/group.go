package model


type Group struct {
	Id int `db:"group_id" json:"group_id"`
	Name string `db:"group_name" json:"group_name"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}