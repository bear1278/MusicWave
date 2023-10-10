package repository

import (
	"database/sql"
	"fmt"
	music "github.com/bear1278/MusicWave"
)

type AuthMysql struct {
	DB *sql.DB
}

func NewAuthMysql(DB *sql.DB) *AuthMysql {
	return &AuthMysql{DB: DB}
}

func (a *AuthMysql) CreateUser(user music.User) (int64, error) {
	var id int64
	query := fmt.Sprintf("INSERT INTO %s (name,username,email,password) values (?,?,?,?)", usersTable)
	dbResult, err := a.DB.Exec(query, user.Name, user.UserName, user.Email, user.Password)
	if err != nil {
		return 0, err
	}
	id, err = dbResult.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (a *AuthMysql) GetUser(username, password string) (music.User, error) {
	var user music.User
	query := fmt.Sprintf("SELECT id,name,email from %s where username=? and password=?", usersTable)
	rows, err := a.DB.Query(query, username, password)
	if err != nil {
		return user, err
	}
	rows.Next()
	err = rows.Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return user, err
	}
	user.UserName = username
	user.Password = password
	return user, nil
}
