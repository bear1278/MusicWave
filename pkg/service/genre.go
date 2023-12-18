package service

import (
	music "github.com/bear1278/MusicWave"
	"github.com/bear1278/MusicWave/pkg/repository"
)

type GenreServiceImpl struct {
	repo repository.GenreRepo
}

func NewGenreServiceImpl(repo repository.GenreRepo) *GenreServiceImpl {
	return &GenreServiceImpl{repo: repo}
}

func (g *GenreServiceImpl) AddGenreToArtist(artistID string, genres []music.Genre) error {
	return g.repo.InsertGenre(artistID, genres)
}

func (g *GenreServiceImpl) GerUserGenres(userId int64) ([]music.Genre, error) {
	return g.repo.SelectUserGenre(userId)
}
