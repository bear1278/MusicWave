package repository

import (
	"database/sql"
	"errors"
	"fmt"
	music "github.com/bear1278/MusicWave"
	"time"
)

type AuthMysql struct {
	DB *sql.DB
}

func NewAuthMysql(DB *sql.DB) *AuthMysql {
	return &AuthMysql{DB: DB}
}

func (a *AuthMysql) CreateUser(user music.User) (int64, error) {
	var id int64
	transaction, err := a.DB.Begin()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (name,username,email,password,picture,create_date,modified_date) values (?,?,?,?,?,?,?)", usersTable)
	dbResult, err := transaction.Exec(query, user.Name, user.UserName, user.Email, user.Password, "/static/images/user.svg", time.Now(), time.Now())
	if err != nil {
		transaction.Rollback()
		return 0, err
	}
	id, err = dbResult.LastInsertId()
	if err != nil {
		transaction.Rollback()
		return 0, err
	}
	query = fmt.Sprintf("INSERT INTO %s (name,duration,cover,author,release_date,type,modified_date) values (?,?,?,?,?,?,?)", playlistTable)
	_, err = transaction.Exec(query, "favorites", 0, "/static/images/heart.svg", id, time.Now(), "private", time.Now())
	if err != nil {
		transaction.Rollback()
		return 0, err
	}
	return id, transaction.Commit()
}

func (a *AuthMysql) GetUser(username, password string) (music.User, error) {
	var user music.User
	query := fmt.Sprintf("SELECT id,name,email,picture from %s where username=? and password=?", usersTable)
	rows, err := a.DB.Query(query, username, password)
	if err != nil {
		return user, err
	}
	if rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Picture)
	} else {
		return user, errors.New("wrong password or username")
	}
	if err != nil {
		return user, err
	}
	user.UserName = username
	user.Password = password
	return user, err
}

func (a *AuthMysql) GetAllGenres() ([]music.Genre, error) {
	var genres []music.Genre
	query := fmt.Sprintf("SELECT * FROM %s", genresTable)
	rows, err := a.DB.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var genre music.Genre
		err := rows.Scan(&genre.Id, &genre.Name)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	return genres, nil
}

func (a *AuthMysql) InsertUserGenre(genres []music.Genre, userId int64) error {
	transaction, err := a.DB.Begin()
	if err != nil {
		return err
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE ID_USER=?", userGenreTable)
	_, err = transaction.Exec(query, userId)
	if err != nil {
		return err
		transaction.Rollback()
	}
	query = fmt.Sprintf("INSERT INTO %s values (?,?)", userGenreTable)
	for _, genre := range genres {
		_, err := transaction.Exec(query, userId, genre.Id)
		if err != nil {
			transaction.Rollback()
			return err
		}
	}
	return transaction.Commit()
}

func (a *AuthMysql) SelectUsersEmail(username, email string) (music.User, error) {
	var user music.User
	query := fmt.Sprintf("SELECT id,name,email from %s where username=? and email=?", usersTable)
	rows, err := a.DB.Query(query, username, email)
	if err != nil {
		return user, err
	}
	if rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Email)
	} else {
		return user, errors.New("wrong username")
	}
	if err != nil {
		return user, err
	}
	user.UserName = username
	return user, err
}

func (a *AuthMysql) SetNewPassword(id int64, password string) error {
	query := fmt.Sprintf("UPDATE %s SET password=?,modified_date=? WHERE id=?", usersTable)
	_, err := a.DB.Exec(query, password, time.Now(), id)
	return err
}
