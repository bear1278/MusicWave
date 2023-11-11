package repository

import (
	"database/sql"
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
	query := fmt.Sprintf("INSERT INTO %s (name,username,email,password) values (?,?,?,?)", usersTable)
	dbResult, err := transaction.Exec(query, user.Name, user.UserName, user.Email, user.Password)
	if err != nil {
		err = transaction.Rollback()
		return 0, err
	}
	id, err = dbResult.LastInsertId()
	if err != nil {
		err = transaction.Rollback()
		return 0, err
	}
	query = fmt.Sprintf("INSERT INTO %s (name,duration,cover,author,release_date) values ('favorites',?,?,?,?)", playlistTable)
	_, err = transaction.Exec(query, 0, "nil", id, time.Now())
	if err != nil {
		err = transaction.Rollback()
		return 0, err
	}
	return id, transaction.Commit()
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

func (a *AuthMysql) GetAllGenres() ([]music.Genre, error) {
	var genres []music.Genre
	query := fmt.Sprintf("SELECT * FROM %s", genresTable)
	rows, err := a.DB.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var genre music.Genre
		err := rows.Scan(&genre.Id, &genre.Name, &genre.Description)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	return genres, nil
}

func (a *AuthMysql) InsertUserGenre(genres []music.Genre, userId int64) error {
	query := fmt.Sprintf("INSERT INTO %s values (?,?)", userGenreTable)
	for _, genre := range genres {
		_, err := a.DB.Exec(query, userId, genre.Id)
		if err != nil {
			return err
		}
	}
	return nil
}
