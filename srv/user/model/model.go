package model

import (
	"database/sql"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

// IsEmailExist
func (mdl *UserModel) IsEmailExist(email string) (bool, error) {
	sqlStr := `select 1 from users where email = ? limit 1`
	stmt, err := mdl.DB.Prepare(sqlStr)
	if err != nil {
		return false, err
	}
	defer stmt.Close()
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
	ImageUrl     string
}

// GetUserByEmail
func (mdl *UserModel) GetUserSecret(email string) (*SecretUser, error) {
	sqlStr := `select ID, Email, Password, PasswordSalt,
Mobile, Nickname, ImageUrl
		from users 
		where Email = ? and isDeleted = 0;
	`
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
	var sUser SecretUser
	for rows.Next() {
		if err := rows.Scan(&sUser.UserId,
			&sUser.Email,
			&sUser.Password,
			&sUser.PasswordSalt,
			&sUser.Mobile,
			&sUser.NickName,
			&sUser.ImageUrl); err != nil {
			return nil, err
		}
	}
	return &sUser, nil
}

type InserUserObj struct {
	Email        string
	Password     string
	PasswordSalt string
	Mobile       string
	Nickname     string
	ImageUrl     string
}

// InserUser
func (mdl *UserModel) InsertUser(user *InserUserObj) (int64, error) {
	sqlStr := `insert into users (Email, Password, PasswordSalt, Mobile, NickName,
ImageUrl, CreatedOn) 
		select ?, ?, ?, ?, ?, ?, ? from dual
		where not exists (
			select 1 from users where Email = ? 
			and IsDeleted = 0
		)
	`
	stmt, err := mdl.DB.Prepare(sqlStr)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	if rst, err := stmt.Exec(
		user.Email,
		user.Password,
		user.PasswordSalt,
		user.Mobile,
		user.Nickname,
		user.ImageUrl,
		time.Now().Format("2006-01-02 15:04:05"),
		user.Email,
	); err == nil {
		return rst.LastInsertId()
	} else {
		return 0, err
	}
}

// UpdateLastLoginTime
func (mdl *UserModel) UpdateLastLoginDate(userId int64) error {
	sqlStr := `update users set LastLoginDate = ? where ID = ?`
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
