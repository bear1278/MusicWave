package service

import (
	"errors"
	music "github.com/bear1278/MusicWave"
	"github.com/bear1278/MusicWave/pkg/repository"
)

type AlbumServiceImpl struct {
	repo repository.AlbumRepo
}

func NewAlbumServiceImpl(repo repository.AlbumRepo) *AlbumServiceImpl {
	return &AlbumServiceImpl{repo: repo}
}

func (a *AlbumServiceImpl) GetAlbumById(idAlbum string) (music.Album, error) {
	return a.repo.GetAlbumById(idAlbum)
}

func (a *AlbumServiceImpl) GetAllAlbumsForUser(idUser int64) ([]music.Album, error) {
	return a.repo.SelectAllAlbumsForUser(idUser)
}

func (a *AlbumServiceImpl) GetAlbumsByArtist(idArtist string) ([]music.Album, error) {
	return a.repo.SelectAlbumsByArtist(idArtist)
}

func (a *AlbumServiceImpl) ExcludeAlbum(idUser int64, idAlbum string) error {
	return a.repo.DeleteAlbumFromFav(idAlbum, idUser)
}

func (a *AlbumServiceImpl) AddAlbumToDB(album music.Album) error {
	tx, err := a.repo.StartTransaction()
	if err != nil {
		return err
	}
	err = a.repo.InsertAlbumToDB(tx, album)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (a *AlbumServiceImpl) AddAlbumToFav(idUser int64, album music.Album) error {
	isExist, err := a.repo.IsAlbumExistInUserLibrary(album.Id, idUser)
	if err != nil {
		return err
	}
	if isExist {
		return errors.New("you already have this album in fav")
	}
	return a.repo.AddAlbumToFav(idUser, album)
}

func (a *AlbumServiceImpl) GetAlbumsRec(userId int64) ([]music.Album, error) {
	albums, err := a.repo.SelectAlbumsByGenre(userId)
	if err != nil {
		return nil, err
	}
	for k, _ := range albums {
		isAllTracks, err := a.repo.IsAllTracksInDB(albums[k].Id, albums[k].TotalTracks)
		if err != nil {
			return nil, err
		}
		if isAllTracks {
			albums[k].Duration, err = a.repo.SelectAlbumDuration(albums[k].Id)
			if err != nil {
				return nil, err
			}
		} else {
			albums[k].Duration = 0
		}
	}
	return albums, err
}
