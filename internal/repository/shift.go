package repository

import (
	"database/sql"

	"shiftmanager/internal/core/db"
	"shiftmanager/internal/model"
)


type ShiftRepository interface {
	Get(s *model.Shift) ([]model.Shift, error)
	GetOne(s *model.Shift) (model.Shift, error)
	Insert(s *model.Shift, tx *sql.Tx) error
	Update(s *model.Shift, tx *sql.Tx) error
	Delete(s *model.Shift, tx *sql.Tx) error
}


type shiftRepository struct {
	db *sql.DB
}

func NewShiftRepository() ShiftRepository {
	db := db.GetDB()
	return &shiftRepository{db}
}


func (rep *shiftRepository) Get(s *model.Shift) ([]model.Shift, error) {
	where, binds := db.BuildWhereClause(s)
	query := 
	`SELECT
		group_id
		,year
		,month
		,store_holiday
		,shift_data
		,created_at
		,updated_at
	 FROM shift ` + where
	rows, err := rep.db.Query(query, binds...)
	defer rows.Close()

	if err != nil {
		return []model.Shift{}, err
	}

	ret := []model.Shift{}
	for rows.Next() {
		s := model.Shift{}
		err = rows.Scan(
			&s.GroupId,
			&s.Year,
			&s.Month,
			&s.StoreHoliday,
			&s.Data,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return []model.Shift{}, err
		}
		ret = append(ret, s)
	}

	return ret, nil
}


func (rep *shiftRepository) GetOne(s *model.Shift) (model.Shift, error) {
	var ret model.Shift
	where, binds := db.BuildWhereClause(s)
	query := 
	`SELECT
		group_id
		,year
		,month
		,store_holiday
		,shift_data
		,created_at
		,updated_at
	 FROM shift ` + where

	err := rep.db.QueryRow(query, binds...).Scan(
		&ret.GroupId,
		&ret.Year,
		&ret.Month,
		&ret.StoreHoliday,
		&ret.Data,
		&ret.CreatedAt,
		&ret.UpdatedAt,
	)

	return ret, err
}


func (rep *shiftRepository) Insert(s *model.Shift, tx *sql.Tx) error {
	cmd := 
	`INSERT INTO shift (
		group_id
		,year
		,month
		,store_holiday
		,shift_data
	 ) VALUES(?,?,?,?,?)`

	binds := []interface{}{
		s.GroupId,
		s.Year,
		s.Month,
		s.StoreHoliday,
		s.Data,
	}

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}


func (rep *shiftRepository) Update(s *model.Shift, tx *sql.Tx) error {
	cmd := 
	`UPDATE shift
	 SET group_id = ?
		,year = ?
		,month = ?
		,store_holiday = ?
		,shift_data = ?
	 WHERE `
	binds := []interface{}{
		s.GroupId,
		s.Year,
		s.Month,
		s.StoreHoliday,
		s.Data,
	}
	
	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}


func (rep *shiftRepository) Delete(s *model.Shift, tx *sql.Tx) error {
	where, binds := db.BuildWhereClause(s)
	cmd := "DELETE FROM shift " + where

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}