package repository

import (
	"database/sql"

	"shiftmanager/internal/core/db"
	"shiftmanager/internal/model"
)


type AccountRepository interface {
	Get(a *model.Account) ([]model.Account, error)
	GetOne(a *model.Account) (model.Account, error)
	Insert(a *model.Account, tx *sql.Tx) error
	Update(a *model.Account, tx *sql.Tx) error
	Delete(a *model.Account, tx *sql.Tx) error
}


type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository() AccountRepository {
	db := db.GetDB()
	return &accountRepository{db}
}


func (rep *accountRepository) Get(a *model.Account) ([]model.Account, error) {
	where, binds := db.BuildWhereClause(a)
	query := 
	`SELECT
		account_id
		,account_name
		,account_password
		,created_at
		,updated_at
	 FROM account ` + where
	rows, err := rep.db.Query(query, binds...)
	defer rows.Close()

	if err != nil {
		return []model.Account{}, err
	}

	ret := []model.Account{}
	for rows.Next() {
		a := model.Account{}
		err = rows.Scan(
			&a.Id,
			&a.Name,
			&a.Password,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return []model.Account{}, err
		}
		ret = append(ret, a)
	}

	return ret, nil
}


func (rep *accountRepository) GetOne(a *model.Account) (model.Account, error) {
	var ret model.Account
	where, binds := db.BuildWhereClause(a)
	query := 
	`SELECT
		account_id
		,account_name
		,account_password
		,created_at
		,updated_at
	 FROM account ` + where

	err := rep.db.QueryRow(query, binds...).Scan(
		&ret.Id,
		&ret.Name,
		&ret.Password,
		&ret.CreatedAt,
		&ret.UpdatedAt,
	)

	return ret, err
}


func (rep *accountRepository) Insert(a *model.Account, tx *sql.Tx) error {
	cmd := 
	`INSERT INTO account (
		account_name
		,account_password
	 ) VALUES(?,?)`

	binds := []interface{}{
		a.Name,
		a.Password,
	}

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}


func (rep *accountRepository) Update(a *model.Account, tx *sql.Tx) error {
	cmd := 
	`UPDATE account
	 SET account_name = ?
		,account_password = ?
	 WHERE account_id = ?`
	binds := []interface{}{
		a.Name,
		a.Password,
		a.Id,
	}
	
	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}


func (rep *accountRepository) Delete(a *model.Account, tx *sql.Tx) error {
	where, binds := db.BuildWhereClause(a)
	cmd := "DELETE FROM account " + where

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}