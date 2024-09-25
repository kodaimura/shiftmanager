package repository

import (
	"database/sql"

	"shiftmanager/internal/core/db"
	"shiftmanager/internal/model"
)


type GroupRepository interface {
	Get(g *model.Group) ([]model.Group, error)
	GetOne(g *model.Group) (model.Group, error)
	Insert(g *model.Group, tx *sql.Tx) (int, error)
	Update(g *model.Group, tx *sql.Tx) error
	Delete(g *model.Group, tx *sql.Tx) error
}


type groupRepository struct {
	db *sql.DB
}

func NewGroupRepository() GroupRepository {
	db := db.GetDB()
	return &groupRepository{db}
}


func (rep *groupRepository) Get(g *model.Group) ([]model.Group, error) {
	where, binds := db.BuildWhereClause(g)
	query := 
	`SELECT
		group_id
		,group_name
		,created_at
		,updated_at
	 FROM group ` + where
	rows, err := rep.db.Query(query, binds...)
	defer rows.Close()

	if err != nil {
		return []model.Group{}, err
	}

	ret := []model.Group{}
	for rows.Next() {
		g := model.Group{}
		err = rows.Scan(
			&g.Id,
			&g.Name,
			&g.CreatedAt,
			&g.UpdatedAt,
		)
		if err != nil {
			return []model.Group{}, err
		}
		ret = append(ret, g)
	}

	return ret, nil
}


func (rep *groupRepository) GetOne(g *model.Group) (model.Group, error) {
	var ret model.Group
	where, binds := db.BuildWhereClause(g)
	query := 
	`SELECT
		group_id
		,group_name
		,created_at
		,updated_at
	 FROM group ` + where

	err := rep.db.QueryRow(query, binds...).Scan(
		&ret.Id,
		&ret.Name,
		&ret.CreatedAt,
		&ret.UpdatedAt,
	)

	return ret, err
}


func (rep *groupRepository) Insert(g *model.Group, tx *sql.Tx) (int, error) {
	cmd := 
	`INSERT INTO group (
		group_name
	 ) VALUES(?)
	 RETURNING group_id`

	binds := []interface{}{
		g.Name,
	}

	var groupId int
	var err error
	if tx != nil {
		err = tx.QueryRow(cmd, binds...).Scan(&groupId)
	} else {
		err = rep.db.QueryRow(cmd, binds...).Scan(&groupId)
	}

	return groupId, err
}


func (rep *groupRepository) Update(g *model.Group, tx *sql.Tx) error {
	cmd := 
	`UPDATE group
	 SET group_name = ?
	 WHERE group_id = ?`
	binds := []interface{}{
		g.Name,
		g.Id,
	}

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }

	return err
}


func (rep *groupRepository) Delete(g *model.Group, tx *sql.Tx) error {
	where, binds := db.BuildWhereClause(g)
	cmd := "DELETE FROM group " + where

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }

	return err
}