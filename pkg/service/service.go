package service

import (
	music "github.com/bear1278/MusicWave"
	"github.com/bear1278/MusicWave/pkg/repository"
)

type Authorization interface {
	CreateUser(user music.User) (int64, error)
	GenerateToken(username, password string) (string, error)
}

type Service struct {
	Authorization
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
	}
}
