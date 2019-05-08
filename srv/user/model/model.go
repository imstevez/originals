package model

import (
	"database/sql"
	"time"

	"shendu.com/errors"

	"github.com/go-redis/redis"
)

type UserModel struct {
	DB *sql.DB
}

const (
	CancelTokenField = `CancelToken`
	FreshTokenField  = `FreshToken`
	RefreshSeconds   = 30
)

var ErrKeyNotExist = errors.New("key is not exist")

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

// CancelToken
func (mdl *UserModel) CancelToken(token string, expiredAt time.Time) error {
	key := CancelTokenField + ":" + token
	expiration := time.Until(expiredAt.Add(RefreshSeconds * time.Second))
	if err := mdl.Redis.Set(key, token, expiration).Err(); err != nil {
		return err
	}
	return nil
}

// IsTokenCanceled
func (mdl *UserModel) IsTokenCanceled(token string) (bool, error) {
	key := CancelTokenField + ":" + token
	val, err := mdl.Redis.Exists(key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, err
}

// GetFreshToken
func (mdl *UserModel) GetFreshToken(token string) (string, error) {
	key := FreshTokenField + ":" + token
	val, err := mdl.Redis.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrKeyNotExist
		}
		return "", err
	}
	return val, err
}

// SetFreshToken
func (mdl *UserModel) SetFreshToken(oldToken, newToken string) error {
	key := FreshTokenField + ":" + oldToken
	if err := mdl.Redis.Set(key, newToken, RefreshSeconds*time.Second).Err(); err != nil {
		return err
	}
	return nil
}
