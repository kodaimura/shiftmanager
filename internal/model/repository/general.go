package repository

import (
	"log"
	"database/sql"

	"shiftmanager/internal/model/db"
	"shiftmanager/internal/model/entity"
)


type GeneralRepository interface {
    SelectByClass(class string) ([]entity.General, error)
}


type generalRepository struct {
	db *sql.DB
}


func NewGeneralRepository() GeneralRepository {
	db := db.GetDB()
	return &generalRepository{db}
}

func (gr *generalRepository)SelectByClass(class string) ([]entity.General, error){
	var ls []entity.General

	rows, err := gr.db.Query(
		`SELECT CLASS, KEY1, KEY2, VALUE1, VALUE2, REMARKS 
		FROM GENERALS WHERE CLASS = ?`, class,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		g := entity.General{}
		err = rows.Scan(&g.Class, &g.Key1, &g.Key2, &g.Value1, &g.Value2, &g.Remarks)
		if err != nil {
			break
		}
		ls = append(ls, g)
	}

	return ls, err
}


func insert(db *sql.DB, class, key1, key2, value1, value2, remarks string) {
	_, err := db.Exec(
		`INSERT INTO GENERALS (CLASS, KEY1, KEY2, VALUE1, VALUE2, REMARKS) 
		VALUES(?,?,?,?,?,?)`, class, key1, key2, value1, value2, remarks,
	)

	if err != nil {
		log.Panic(err)
	}
}


//データ登録
func init() {
	db := db.GetDB()

	cmd := `
		CREATE TABLE IF NOT EXISTS GENERALS (
			CLASS VARCHAR(50),
			KEY1 VARCHAR(30),
			KEY2 VARCHAR(30),
			VALUE1 VARCHAR(50),
			VALUE2 VARCHAR(50),
			REMARKS VARCHAR(10),
			CREATE_AT TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
			UPDATE_AT TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
		);

		CREATE TRIGGER IF NOT EXISTS trigger_generals_updated_at AFTER UPDATE ON GENERALS
		BEGIN
    		UPDATE GROUPS 
    			SET UPDATE_AT = DATETIME('now', 'localtime') 
    			WHERE rowid == NEW.rowid;
		END;`

	db.Exec(cmd)
	_, err := db.Exec(`DELETE FROM GENERALS`)

	if err != nil{
		log.Panic(err)
	}

	insert(db, "role", "1", "", "キッチン", "", "")
	insert(db, "role", "2", "", "ホール", "", "")
	insert(db, "role", "3", "", "キッチン・ホール", "", "")
}