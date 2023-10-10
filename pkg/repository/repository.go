package repository

import (
	"database/sql"
	music "github.com/bear1278/MusicWave"
)

type Authorization interface {
	CreateUser(user music.User) (int64, error)
	GetUser(username, password string) (music.User, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthMysql(db),
	}
}
