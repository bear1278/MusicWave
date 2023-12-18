package repository

import (
	"database/sql"
	"fmt"
	music "github.com/bear1278/MusicWave"
	"time"
)

type UserMysql struct {
	db *sql.DB
}

func NewUserMysql(db *sql.DB) *UserMysql {
	return &UserMysql{db: db}
}

func (u *UserMysql) UpdateUsername(userId int64, username string) error {
	query := fmt.Sprintf("UPDATE %s SET username=?,modified_date=? WHERE id=?", usersTable)
	_, err := u.db.Exec(query, username, time.Now(), userId)
	return err
}

func (u *UserMysql) UpdatePassword(userId int64, password string) error {
	query := fmt.Sprintf("UPDATE %s SET password=?,modified_date=? WHERE id=?", usersTable)
	_, err := u.db.Exec(query, password, time.Now(), userId)
	return err
}

func (u *UserMysql) SelectPassword(userId int64) (string, error) {
	var password string
	query := fmt.Sprintf("SELECT password FROM %s WHERE id=?", usersTable)
	rows, err := u.db.Query(query, userId)
	if err != nil {
		return "", err
	}
	for rows.Next() {
		err = rows.Scan(&password)
		if err != nil {
			return "", err
		}
	}
	return password, err
}

func (u *UserMysql) UpdatePicture(userId int64, picture string) error {
	query := fmt.Sprintf("UPDATE %s SET picture=?,modified_date=? WHERE id=?", usersTable)
	_, err := u.db.Exec(query, picture, time.Now(), userId)
	return err
}

func (u *UserMysql) UpdateEmail(userId int64, email string) error {
	query := fmt.Sprintf("UPDATE %s SET email=?,modified_date=? WHERE id=?", usersTable)
	_, err := u.db.Exec(query, email, time.Now(), userId)
	return err
}

func (u *UserMysql) SelectAllUsers() ([]music.User, error) {
	var user music.User
	var users []music.User
	var createDate, modDate string
	query := fmt.Sprintf("SELECT id,name,username,email,picture,create_date, modified_date FROM %s where username!='admin'", usersTable)
	rows, err := u.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.UserName, &user.Email, &user.Picture, &createDate, &modDate)
		if err != nil {
			return nil, err
		}
		user.CreateDate, err = time.Parse(time.DateTime, createDate)
		if err != nil {
			return nil, err
		}
		user.ModifiedDate, err = time.Parse(time.DateTime, modDate)
		if err != nil {
			return nil, err
		}
		if user.Picture == "" || user.Picture == " " {
			user.Picture = "/static/images/user.svg"
		}
		users = append(users, user)
		user = music.User{}
	}
	return users, err
}

func (u *UserMysql) DeleteUser(idUser int64, reason string) error {
	var username string
	transaction, err := u.db.Begin()
	if err != nil {
		return err
	}
	query := fmt.Sprintf("SELECT username FROM %s WHERE id=?", usersTable)
	rows, err := transaction.Query(query, idUser)
	if err != nil {
		transaction.Rollback()
		return err
	}
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			transaction.Rollback()
			return err
		}
	}
	query = fmt.Sprintf("INSERT INTO %s (name,type,reason,deleted_date) values (?,?,?,?)", adminHistory)
	_, err = transaction.Exec(query, username, "user", reason, time.Now())
	if err != nil {
		transaction.Rollback()
		return err
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE author=?", playlistTable)
	_, err = transaction.Exec(query, idUser)
	if err != nil {
		transaction.Rollback()
		return err
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE id=?", usersTable)
	_, err = transaction.Exec(query, idUser)
	if err != nil {
		transaction.Rollback()
		return err
	}
	return transaction.Commit()
}

func (u *UserMysql) SelectUserById(idUser int64) (music.User, error) {
	var user music.User
	query := fmt.Sprintf("SELECT id,username,email,picture from %s where id=?", usersTable)
	rows, err := u.db.Query(query, idUser)
	if err != nil {
		return music.User{}, err
	}
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.UserName, &user.Email, &user.Picture)
		if err != nil {
			return music.User{}, err
		}
		if user.Picture == "" {
			user.Picture = "/static/images/user.svg"
		}
	}
	return user, err
}
