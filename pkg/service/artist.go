package service

import (
	"errors"
	music "github.com/bear1278/MusicWave"
	"github.com/bear1278/MusicWave/pkg/repository"
)

type ArtistServiceImpl struct {
	repo repository.ArtistRepo
}

func (a *ArtistServiceImpl) GetAllAlbumsForUser(idUser int64) ([]music.Artist, error) {
	return a.repo.SelectAllArtistForUser(idUser)
}

func NewArtistServiceImpl(repo repository.ArtistRepo) *ArtistServiceImpl {
	return &ArtistServiceImpl{repo: repo}
}

func (a *ArtistServiceImpl) GetArtistById(idArtist string) (music.Artist, error) {
	return a.repo.SelectArtistById(idArtist)
}

func (a *ArtistServiceImpl) ExcludeArtist(idUser int64, idArtist string) error {
	return a.repo.DeleteArtistFromFav(idArtist, idUser)
}

func (a *ArtistServiceImpl) AddArtistToFav(idUser int64, artist music.Artist) error {
	isExist, err := a.repo.IsArtistExistInUserLibrary(artist.Id, idUser)
	if err != nil {
		return err
	}
	if isExist {
		return errors.New("you already have this artist in your library")
	} else {
		return a.repo.AddArtistToFav(idUser, artist)
	}
}

func (a *ArtistServiceImpl) AddArtistToDB(artist music.Artist) error {
	tx, err := a.repo.StartTransaction()
	if err != nil {
		return err
	}
	err = a.repo.InsertArtistToDB(tx, artist)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (a *ArtistServiceImpl) GetAllArtists() ([]music.Artist, error) {
	return a.repo.SelectAllArtists()
}

func (a *ArtistServiceImpl) DeleteArtist(artistID, reason string) error {
	return a.repo.DeleteArtist(artistID, reason)
}

func (a *ArtistServiceImpl) GetArtistsRec(userId int64) ([]music.Artist, error) {
	return a.repo.SelectArtistsByGenre(userId)
}

func (a *ArtistServiceImpl) GetMostPopularTrackForUser(userId int64) ([]music.Artist, error) {
	return a.repo.SelectMostPopularArtistForUser(userId)
}
