package repository

import (
	"log"
	"database/sql"

	"shiftmanager/internal/model/db"
	"shiftmanager/internal/model/entity"
)


func init(){
	db := db.GetDB()

	cmd := `
		CREATE TABLE IF NOT EXISTS SHIFTS (
			GID INTEGER,
			YEAR VARCHAR(4),
			MONTH VARCHAR(2),
			STORE_HOLIDAY VARCHAR(100),
			SHIFT VARCHAR(500),
			CREATE_AT TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
			UPDATE_AT TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
			PRIMARY KEY(GID, YEAR, MONTH)
		);

		CREATE TRIGGER IF NOT EXISTS trigger_shifts_updated_at AFTER UPDATE ON SHIFTS
		BEGIN
    		UPDATE SHIFTS
    			SET UPDATE_AT = DATETIME('now', 'localtime') 
    			WHERE rowid == NEW.rowid;
		END;`

	_, err := db.Exec(cmd)

	if err != nil {
		log.Panic(err)
	}
}


type ShiftRepository interface {
	SelectByGIdYearMonth(gid int, year string, month string) (entity.Shift, error)
    Upsert(shift entity.Shift) error
}


type shiftRepository struct {
	db *sql.DB
}


func NewShiftRepository() ShiftRepository {
	db := db.GetDB()
	return &shiftRepository{db}
}


func (sr *shiftRepository)SelectByGIdYearMonth(gid int, year string, month string) (entity.Shift, error){
	var s entity.Shift
	err := sr.db.QueryRow(
		`SELECT 
			GID, YEAR, MONTH, STORE_HOLIDAY,
		 	SHIFT, CREATE_AT, UPDATE_AT 
		 FROM SHIFTS WHERE GID = ? AND YEAR = ? AND MONTH = ?`, 
		gid, year, month,
	).Scan(
		&s.GId, &s.Year, &s.Month, &s.StoreHoliday,
		&s.Shift, &s.CreateAt, &s.UpdateAt,
	)

	return s, err
}


func (sr *shiftRepository)Upsert(shift entity.Shift) error {
	_, err := sr.db.Exec(
		`INSERT OR REPLACE INTO SHIFTS 
			(GID, YEAR, MONTH, STORE_HOLIDAY, SHIFT)  
		 VALUES(
			?,?,?,?,?
		)`,
		shift.GId, shift.Year, shift.Month, shift.StoreHoliday, shift.Shift,
	)

	return err
}
