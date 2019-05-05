package model

import (
	"database/sql"
	"originals/srv/user/proto"
	"time"
)

type UserSrvModel struct {
	DB *sql.DB
}

// GetUserByEmail
func (mdl *UserSrvModel) GetUserByEmail(email string) (*proto.User, error) {
	sqlStr := `select ID, Email, Password, PasswordSalt, Mobile, Nickname, ImageUrl 
from users where Email = ? and isDeleted = 0;`
	stmt, err := mdl.DB.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var user proto.User
	for rows.Next() {
		if err := rows.Scan(&user.Id,
			&user.Email,
			&user.Password,
			&user.PasswordSalt,
			&user.Mobile,
			&user.Nickname,
			&user.ImageUrl); err != nil {
			return nil, err
		}
	}
	return &user, nil
}

// CountUserEmail
func (mdl *UserSrvModel) CountUserEmail(email string) (int, error) {
	sqlStr := `select count(Email) from users where email = ?;`
	stmt, err := mdl.DB.Prepare(sqlStr)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(email)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
	}
	return count, nil
}

// InserUser
func (mdl *UserSrvModel) InsertUser(user *proto.User) (int64, error) {
	sqlStr := `insert into users (Email, Password, PasswordSalt, NickName,
 ImageUrl, CreatedOn) Values (?, ?, ?, ?, ?, ?)`
	stmt, err := mdl.DB.Prepare(sqlStr)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	if rst, err := stmt.Exec(
		user.Email,
		user.Password,
		user.PasswordSalt,
		user.Nickname,
		user.ImageUrl,
		time.Now().Format("2006-01-02 15:04:05"),
	); err == nil {
		return rst.LastInsertId()
	} else {
		return 0, err
	}
}

// UpdateLastLoginTime
func (mdl *UserSrvModel) UpdateLastLoginTime(userId int64) error {
	sqlStr := `update users set LastLoginTime = ? where ID = ?`
	stmt, err := mdl.DB.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(userId,
		time.Now().Format("2006-01-02 15:04:05")); err != nil {
		return err
	}
	return nil
}
