package repository

import (
	"database/sql"

	"shiftmanager/internal/core/db"
	"shiftmanager/internal/model"
)


type AvailableWorkDaysRepository interface {
	Get(awd *model.AvailableWorkDays) ([]model.AvailableWorkDays, error)
	GetOne(awd *model.AvailableWorkDays) (model.AvailableWorkDays, error)
	Insert(awd *model.AvailableWorkDays, tx *sql.Tx) error
	Update(awd *model.AvailableWorkDays, tx *sql.Tx) error
	Delete(awd *model.AvailableWorkDays, tx *sql.Tx) error
}


type availableWorkDaysRepository struct {
	db *sql.DB
}

func NewAvailableWorkDaysRepository() AvailableWorkDaysRepository {
	db := db.GetDB()
	return &availableWorkDaysRepository{db}
}


func (rep *availableWorkDaysRepository) Get(awd *model.AvailableWorkDays) ([]model.AvailableWorkDays, error) {
	where, binds := db.BuildWhereClause(awd)
	query := 
	`SELECT
		account_id
		,year
		,month
		,available_work_days
		,memo
		,created_at
		,updated_at
	 FROM available_work_days ` + where
	rows, err := rep.db.Query(query, binds...)
	defer rows.Close()

	if err != nil {
		return []model.AvailableWorkDays{}, err
	}

	ret := []model.AvailableWorkDays{}
	for rows.Next() {
		awd := model.AvailableWorkDays{}
		err = rows.Scan(
			&awd.AccountId,
			&awd.Year,
			&awd.Month,
			&awd.AvailableWorkDays,
			&awd.Memo,
			&awd.CreatedAt,
			&awd.UpdatedAt,
		)
		if err != nil {
			return []model.AvailableWorkDays{}, err
		}
		ret = append(ret, awd)
	}

	return ret, nil
}


func (rep *availableWorkDaysRepository) GetOne(awd *model.AvailableWorkDays) (model.AvailableWorkDays, error) {
	var ret model.AvailableWorkDays
	where, binds := db.BuildWhereClause(awd)
	query := 
	`SELECT
		account_id
		,year
		,month
		,available_work_days
		,memo
		,created_at
		,updated_at
	 FROM available_work_days ` + where

	err := rep.db.QueryRow(query, binds...).Scan(
		&ret.AccountId,
		&ret.Year,
		&ret.Month,
		&ret.AvailableWorkDays,
		&ret.Memo,
		&ret.CreatedAt,
		&ret.UpdatedAt,
	)

	return ret, err
}


func (rep *availableWorkDaysRepository) Insert(awd *model.AvailableWorkDays, tx *sql.Tx) error {
	cmd := 
	`INSERT INTO available_work_days (
		account_id
		,year
		,month
		,available_work_days
		,memo
	 ) VALUES(?,?,?,?,?)`

	binds := []interface{}{
		awd.AccountId,
		awd.Year,
		awd.Month,
		awd.AvailableWorkDays,
		awd.Memo,
	}

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}


func (rep *availableWorkDaysRepository) Update(awd *model.AvailableWorkDays, tx *sql.Tx) error {
	cmd := 
	`UPDATE available_work_days
	 SET account_id = ?
		,year = ?
		,month = ?
		,available_work_days = ?
		,memo = ?
	 WHERE `
	binds := []interface{}{
		awd.AccountId,
		awd.Year,
		awd.Month,
		awd.AvailableWorkDays,
		awd.Memo,
	}
	
	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}


func (rep *availableWorkDaysRepository) Delete(awd *model.AvailableWorkDays, tx *sql.Tx) error {
	where, binds := db.BuildWhereClause(awd)
	cmd := "DELETE FROM available_work_days " + where

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}