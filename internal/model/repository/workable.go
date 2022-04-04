package repository

import (
	"log"
	"database/sql"

	"shiftmanager/internal/dto"
	"shiftmanager/internal/model/db"
	"shiftmanager/internal/model/entity"
)


func init(){
	db := db.GetDB()

	cmd := `
		CREATE TABLE IF NOT EXISTS WORKABLES (
			UID INTEGER,
			YEAR VARCHAR(4),
			MONTH VARCHAR(2),
			WORKABLE_DAYS VARCHAR(100),
			MEMO VARCHAR(100),
			CREATE_AT TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
			UPDATE_AT TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
			PRIMARY KEY(UID, YEAR, MONTH)
		);

		CREATE TRIGGER IF NOT EXISTS trigger_workables_updated_at AFTER UPDATE ON WORKABLES
		BEGIN
    		UPDATE WORKABLES
    			SET UPDATE_AT = DATETIME('now', 'localtime') 
    			WHERE rowid == NEW.rowid;
		END;`

	_, err := db.Exec(cmd)

	if err != nil {
		log.Panic(err)
	}
}


type WorkableRepository interface {
	SelectByUIdYearMonth(uid int, year string, month string) (entity.Workable, error)
    Upsert(workable entity.Workable) error

    GetWorkableExp1ByGId(gid int, year string, month string) ([]dto.WorkableExp1, error)
}


type workableRepository struct {
	db *sql.DB
}


func NewWorkableRepository() WorkableRepository {
	db := db.GetDB()
	return &workableRepository{db}
}


func (wr *workableRepository)SelectByUIdYearMonth(uid int, year string, month string) (entity.Workable, error){
	var w entity.Workable
	err := wr.db.QueryRow(
		`SELECT UID, YEAR, MONTH, WORKABLE_DAYS, MEMO, CREATE_AT, UPDATE_AT 
		 FROM WORKABLES WHERE UID = ? AND YEAR = ? AND MONTH = ?`, 
		uid, year, month,
	).Scan(
		&w.UId, &w.Year, &w.Month, &w.WorkableDays,
		&w.Memo, &w.CreateAt, &w.UpdateAt,
	)

	return w, err
}


func (wr *workableRepository)Upsert(workable entity.Workable) error {
	_, err := wr.db.Exec(
		`REPLACE INTO WORKABLES (UID, YEAR, MONTH, WORKABLE_DAYS, MEMO) 
		 VALUES(?,?,?,?,?)`,
		workable.UId, workable.Year, workable.Month, workable.WorkableDays, workable.Memo,
	)

	return err
}


func (wr *workableRepository)GetWorkableExp1ByGId(gid int, year string, month string) ([]dto.WorkableExp1, error){
	var ls []dto.WorkableExp1
	rows, err := wr.db.Query(
		`SELECT 
			W.UID, W.YEAR, W.MONTH, W.WORKABLE_DAYS, 
			P.ABB_NAME, P.ROLE
		 FROM WORKABLES AS W 
		 LEFT JOIN PROFILES AS P ON W.UID = P.UID
		 WHERE P.GID = ? AND YEAR = ? AND MONTH = ?`, 
		gid, year, month,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		we := dto.WorkableExp1{}
		err = rows.Scan(
			&we.UId, &we.Year, &we.Month, &we.WorkableDays,
			&we.AbbName, &we.Role,
		)
		if err != nil {
			break
		}
		ls = append(ls, we)
	}

	return ls, err
}