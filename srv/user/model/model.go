package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/micro/go-log"
)

type UserModel struct {
	DB *sql.DB
}

func logSQL(sql string, args []interface{}) {
	log.Logf("sql string: \n%s\nsql args: \n%v\n", sql, args)
}

// IsEmailExist
func (mdl *UserModel) IsEmailExist(email string) (bool, error) {
	sqlStr := `select 1 from Users where email = ? limit 1`
	stmt, err := mdl.DB.Prepare(sqlStr)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	logSQL(sqlStr, []interface{}{email})

	rows, err := stmt.Query(email)
	if err != nil {
		return false, nil
	}
	defer rows.Close()
	for rows.Next() {
		return true, nil
	}
	return false, nil
}

type SecretUser struct {
	UserId       int64
	Email        string
	Password     string
	PasswordSalt string
	Mobile       string
	NickName     string
	Avatar       string
}

// GetUserByEmail
func (mdl *UserModel) GetUserSecret(email string) (*SecretUser, error) {
	sqlStr := `select ID, Email, Password, PasswordSalt, Nickname, Avatar
		from Users 
		where Email = ? and isDeleted = 0;
	`
	stmt, err := mdl.DB.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	logSQL(sqlStr, []interface{}{email})

	rows, err := stmt.Query(email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sUser SecretUser
	for rows.Next() {
		if err := rows.Scan(&sUser.UserId,
			&sUser.Email,
			&sUser.Password,
			&sUser.PasswordSalt,
			&sUser.NickName,
			&sUser.Avatar); err != nil {
			return nil, err
		}
	}
	return &sUser, nil
}

type InsertUserObj struct {
	Email        string
	Password     string
	PasswordSalt string
	Nickname     string
	Avatar       string
}

// InsertUser
func (mdl *UserModel) InsertUser(user *InsertUserObj) (int64, error) {
	sqlStr := `insert into Users (Email, Password, PasswordSalt, NickName,
Avatar, CreatedOn) 
		select ?, ?, ?, ?, ?, ? from dual
		where not exists (
			select 1 from Users where Email = ? 
			and IsDeleted = 0
		)
	`
	stmt, err := mdl.DB.Prepare(sqlStr)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	args := []interface{}{
		user.Email,
		user.Password,
		user.PasswordSalt,
		user.Nickname,
		user.Avatar,
		time.Now().Format("2006-01-02 15:04:05"),
		user.Email,
	}

	logSQL(sqlStr, args)

	if rst, err := stmt.Exec(args...); err == nil {
		return rst.LastInsertId()
	} else {
		return 0, err
	}
}

// UpdateLastLoginTime
func (mdl *UserModel) UpdateLastLoginDate(userId int64) error {
	sqlStr := `update Users set LastLoginDate='%s' where ID = '%d';`
	sqlStr = fmt.Sprintf(sqlStr, time.Now().Format("2006-01-02 15:04:05"), userId)
	stmt, err := mdl.DB.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	logSQL(sqlStr, []interface{}{})

	if _, err := stmt.Exec(); err != nil {
		return err
	}
	return nil
}
