package repository

import (
	"database/sql"
	music "github.com/bear1278/MusicWave"
)

type Authorization interface {
	CreateUser(user music.User) (int64, error)
	GetUser(username, password string) (music.User, error)
	GetAllGenres() ([]music.Genre, error)
	InsertUserGenre(genres []music.Genre, userId int64) error
}

type PlaylistRepo interface {
	CreatePlaylist(playlist music.Playlist) (int64, error)
	DeletePlaylist(playlist int64) error
	SelectAllPlaylists(idUser int64) ([]music.Playlist, error)
	SelectById(playlist int64) (music.Playlist, error)
	UpdatePlaylist(playlist music.Playlist) error
	ExcludePlaylist(idPlaylist, idUser int64) error
	AddPlaylist(idPlaylist, idUser int64) error
}

type TrackRepo interface {
	AddTrack(idPlaylist int64, track music.Track) (int64, error)
}

type Repository struct {
	Authorization
	PlaylistRepo
	TrackRepo
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthMysql(db),
		PlaylistRepo:  NewPlaylistMysql(db),
	}
}
