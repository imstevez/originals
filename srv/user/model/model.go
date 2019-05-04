package model

import (
	"database/sql"
	"originals/srv/user/proto"
)

type UserSrvModel struct {
	DB *sql.DB
}

func (mdl *UserSrvModel) GetUserByEmail(email string) ([]*proto.User, error) {
	sqlStr := `select ID, UserName, Email from users where Email = ?;`
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

	var users []*proto.User
	for rows.Next() {
		var user proto.User
		if err := rows.Scan(&user.Id,
			&user.UserName,
			&user.Email); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}
