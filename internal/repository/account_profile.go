package repository

import (
	"database/sql"

	"shiftmanager/internal/core/db"
	"shiftmanager/internal/model"
)


type AccountProfileRepository interface {
	Get(ap *model.AccountProfile) ([]model.AccountProfile, error)
	GetOne(ap *model.AccountProfile) (model.AccountProfile, error)
	Insert(ap *model.AccountProfile, tx *sql.Tx) error
	Update(ap *model.AccountProfile, tx *sql.Tx) error
	Delete(ap *model.AccountProfile, tx *sql.Tx) error
}


type accountProfileRepository struct {
	db *sql.DB
}

func NewAccountProfileRepository() AccountProfileRepository {
	db := db.GetDB()
	return &accountProfileRepository{db}
}


func (rep *accountProfileRepository) Get(ap *model.AccountProfile) ([]model.AccountProfile, error) {
	where, binds := db.BuildWhereClause(ap)
	query := 
	`SELECT
		account_id
		,group_id
		,account_role
		,display_name
		,created_at
		,updated_at
	 FROM account_profile ` + where
	rows, err := rep.db.Query(query, binds...)
	defer rows.Close()

	if err != nil {
		return []model.AccountProfile{}, err
	}

	ret := []model.AccountProfile{}
	for rows.Next() {
		ap := model.AccountProfile{}
		err = rows.Scan(
			&ap.AccountId,
			&ap.GroupId,
			&ap.AccountRole,
			&ap.DisplayName,
			&ap.CreatedAt,
			&ap.UpdatedAt,
		)
		if err != nil {
			return []model.AccountProfile{}, err
		}
		ret = append(ret, ap)
	}

	return ret, nil
}


func (rep *accountProfileRepository) GetOne(ap *model.AccountProfile) (model.AccountProfile, error) {
	var ret model.AccountProfile
	where, binds := db.BuildWhereClause(ap)
	query := 
	`SELECT
		account_id
		,group_id
		,account_role
		,display_name
		,created_at
		,updated_at
	 FROM account_profile ` + where

	err := rep.db.QueryRow(query, binds...).Scan(
		&ret.AccountId,
		&ret.GroupId,
		&ret.AccountRole,
		&ret.DisplayName,
		&ret.CreatedAt,
		&ret.UpdatedAt,
	)

	return ret, err
}


func (rep *accountProfileRepository) Insert(ap *model.AccountProfile, tx *sql.Tx) error {
	cmd := 
	`INSERT INTO account_profile (
		account_id
		,group_id
		,account_role
		,display_name
	 ) VALUES(?,?,?,?)`

	binds := []interface{}{
		ap.AccountId,
		ap.GroupId,
		ap.AccountRole,
		ap.DisplayName,
	}

	var err error
	if tx != nil {
		_, err = tx.Exec(cmd, binds...)
	} else {
		_, err = rep.db.Exec(cmd, binds...)
	}

	return err
}


func (rep *accountProfileRepository) Update(ap *model.AccountProfile, tx *sql.Tx) error {
	cmd := 
	`UPDATE account_profile
	 SET account_id = ?
		,group_id = ?
		,account_role = ?
		,display_name = ?
	 WHERE `
	binds := []interface{}{
		ap.AccountId,
		ap.GroupId,
		ap.AccountRole,
		ap.DisplayName,
	}

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }

	return err
}


func (rep *accountProfileRepository) Delete(ap *model.AccountProfile, tx *sql.Tx) error {
	where, binds := db.BuildWhereClause(ap)
	cmd := "DELETE FROM account_profile " + where

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }

	return err
}