package repository

import (
	"fmt"
	"database/sql"

	"shiftmanager/internal/core/db"
	"shiftmanager/internal/model"
)


type ShiftPreferredRepository interface {
	Get(sp *model.ShiftPreferred) ([]model.ShiftPreferred, error)
	GetOne(sp *model.ShiftPreferred) (model.ShiftPreferred, error)
	Insert(sp *model.ShiftPreferred, tx *sql.Tx) error
	Update(sp *model.ShiftPreferred, tx *sql.Tx) error
	Delete(sp *model.ShiftPreferred, tx *sql.Tx) error
}


type shiftPreferredRepository struct {
	db *sql.DB
}

func NewShiftPreferredRepository() ShiftPreferredRepository {
	db := db.GetDB()
	return &shiftPreferredRepository{db}
}


func (rep *shiftPreferredRepository) Get(sp *model.ShiftPreferred) ([]model.ShiftPreferred, error) {
	where, binds := db.BuildWhereClause(sp)
	query := 
	`SELECT
		account_id
		,year
		,month
		,dates
		,notes
		,created_at
		,updated_at
	 FROM shift_preferred ` + where
	rows, err := rep.db.Query(query, binds...)
	defer rows.Close()

	if err != nil {
		return []model.ShiftPreferred{}, err
	}

	ret := []model.ShiftPreferred{}
	for rows.Next() {
		sp := model.ShiftPreferred{}
		err = rows.Scan(
			&sp.AccountId,
			&sp.Year,
			&sp.Month,
			&sp.Dates,
			&sp.Notes,
			&sp.CreatedAt,
			&sp.UpdatedAt,
		)
		if err != nil {
			return []model.ShiftPreferred{}, err
		}
		ret = append(ret, sp)
	}

	return ret, nil
}


func (rep *shiftPreferredRepository) GetOne(sp *model.ShiftPreferred) (model.ShiftPreferred, error) {
	var ret model.ShiftPreferred
	where, binds := db.BuildWhereClause(sp)
	query := 
	`SELECT
		account_id
		,year
		,month
		,dates
		,notes
		,created_at
		,updated_at
	 FROM shift_preferred ` + where

	fmt.Println(query)

	err := rep.db.QueryRow(query, binds...).Scan(
		&ret.AccountId,
		&ret.Year,
		&ret.Month,
		&ret.Dates,
		&ret.Notes,
		&ret.CreatedAt,
		&ret.UpdatedAt,
	)

	return ret, err
}


func (rep *shiftPreferredRepository) Insert(sp *model.ShiftPreferred, tx *sql.Tx) error {
	cmd := 
	`INSERT INTO shift_preferred (
		account_id
		,year
		,month
		,dates
		,notes
	 ) VALUES(?,?,?,?,?)`

	binds := []interface{}{
		sp.AccountId,
		sp.Year,
		sp.Month,
		sp.Dates,
		sp.Notes,
	}

	var err error
	if tx != nil {
		_, err = tx.Exec(cmd, binds...)
	} else {
		_, err = rep.db.Exec(cmd, binds...)
	}

	return err
}


func (rep *shiftPreferredRepository) Update(sp *model.ShiftPreferred, tx *sql.Tx) error {
	cmd := 
	`UPDATE shift_preferred
	 SET account_id = ?
		,year = ?
		,month = ?
		,dates = ?
		,notes = ?
	 WHERE account_id = ?
	   AND year = ?
	   AND month = ?`
	binds := []interface{}{
		sp.AccountId,
		sp.Year,
		sp.Month,
		sp.Dates,
		sp.Notes,
		sp.AccountId,
		sp.Year,
		sp.Month,
	}

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }

	return err
}


func (rep *shiftPreferredRepository) Delete(sp *model.ShiftPreferred, tx *sql.Tx) error {
	where, binds := db.BuildWhereClause(sp)
	cmd := "DELETE FROM shift_preferred " + where

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }

	return err
}