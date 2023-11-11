package service

import (
	music "github.com/bear1278/MusicWave"
	"github.com/bear1278/MusicWave/pkg/repository"
)

type Authorization interface {
	CreateUser(user music.User) (int64, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int64, error)
	FillHtml() ([]music.Genre, error)
	InsertRecommendation(genres []music.Genre, userId int64) error
}

type PlaylistService interface {
	NewPlaylist(playlist music.Playlist, userId int64) (int64, error)
	DeletePlaylist(playlist int64) error
	GetAllPlaylist(idUser int64) ([]music.Playlist, error)
	GetById(playlist int64) (music.Playlist, error)
	UpdatePlaylist(playlist music.Playlist) error
	ExcludePlaylist(idPlaylist, idUser int64) error
	AddPlaylist(idPlaylist, idUser int64) error
}

type Service struct {
	Authorization
	PlaylistService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization:   NewAuthService(repo.Authorization),
		PlaylistService: NewPlaylistServiceImpl(repo.PlaylistRepo),
	}
}
