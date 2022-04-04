package repository

import (
	"log"
	"database/sql"

	"shiftmanager/internal/model/db"
	"shiftmanager/internal/model/entity"
	"shiftmanager/internal/dto"
)


func init(){
	db := db.GetDB()

	cmd := `
		CREATE TABLE IF NOT EXISTS USERS (
			UID INTEGER PRIMARY KEY AUTOINCREMENT,
			USER_NAME VARCHAR(100) NOT NULL UNIQUE,
			PASSWORD VARCHAR(100) NOT NULL,
			CREATE_AT TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
			UPDATE_AT TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
		);

		CREATE TRIGGER IF NOT EXISTS trigger_users_updated_at AFTER UPDATE ON USERS
		BEGIN
    		UPDATE USERS 
    			SET UPDATE_AT = DATETIME('now', 'localtime') 
    			WHERE rowid == NEW.rowid;
		END;`

	_, err := db.Exec(cmd)

	if err != nil {
		log.Panic(err)
	}
}


type UserRepository interface {
	Select() ([]entity.User, error)
    SelectByUId(uid int) (entity.User, error)
    UpdateByUId(uid int, user entity.User) error
    DeleteByUId(uid int) error

    Signup(sd dto.SignupDto) error
    SelectByUserName(userName string) (entity.User, error)
    UpdatePasswordByUId(uid int, password string) error
    UpdateUserNameByUId(uid int, userName string) error
}


type userRepository struct {
	db *sql.DB
}


func NewUserRepository() UserRepository {
	db := db.GetDB()
	return &userRepository{db}
}


func (ur *userRepository)Select() ([]entity.User, error){
	var ls []entity.User

	rows, err := ur.db.Query(
		`SELECT UID, USER_NAME, CREATE_AT, UPDATE_AT FROM USERS`,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		u := entity.User{}
		err = rows.Scan(&u.UId, &u.UserName, &u.CreateAt, &u.UpdateAt)
		if err != nil {
			break
		}
		ls = append(ls, u)
	}

	return ls, err
}


func (ur *userRepository)SelectByUId(uid int) (entity.User, error){
	var u entity.User
	err := ur.db.QueryRow(
		`SELECT UID, USER_NAME, CREATE_AT, UPDATE_AT FROM USERS WHERE UID = ?`, uid,
	).Scan(
		&u.UId, &u.UserName, &u.CreateAt, &u.UpdateAt,
	)

	return u, err
}


func (ur *userRepository)UpdateByUId(uid int, user entity.User) error {
	_, err := ur.db.Exec(
		`UPDATE USERS SET USER_NAME = ? WHERE UID = ?`,
		user.UserName, uid,
	)
	return err
}


func (ur *userRepository)DeleteByUId(uid int) error {
	_, err := ur.db.Exec(`DELETE FROM USERS WHERE UID = ?`, uid)

	return err
}


func (ur *userRepository)Signup(sd dto.SignupDto) error {
	_, err := ur.db.Exec(
		`INSERT INTO USERS (USER_NAME, PASSWORD) VALUES(?,?)`,
		sd.UserName, sd.Password,
	)

	return err
}


func (ur *userRepository)SelectByUserName(userName string) (entity.User, error){
	var u entity.User
	err := ur.db.QueryRow(
		`SELECT UID, USER_NAME, PASSWORD, CREATE_AT, UPDATE_AT 
		 FROM USERS WHERE USER_NAME = ?`, userName,
	).Scan(
		&u.UId, &u.UserName, &u.Password, &u.CreateAt, &u.UpdateAt,
	)
	return u, err
}


func (ur *userRepository)UpdatePasswordByUId(uid int, password string) error {
	_, err := ur.db.Exec(
		`UPDATE USERS SET PASSWORD = ? WHERE UID = ?`, password, uid,
	)
	return err
}


func (ur *userRepository)UpdateUserNameByUId(uid int, userName string) error {
	_, err := ur.db.Exec(
		`UPDATE USERS SET USER_NAME = ? WHERE UID = ?`, userName, uid,
	)
	return err
}
