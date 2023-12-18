package service

import (
	"errors"
	music "github.com/bear1278/MusicWave"
	"github.com/bear1278/MusicWave/pkg/repository"
)

type TrackServiceImpl struct {
	repo repository.TrackRepo
}

func NewTrackServiceImpl(repo repository.TrackRepo) *TrackServiceImpl {
	return &TrackServiceImpl{repo: repo}
}

func (t *TrackServiceImpl) AddTrack(idPlaylist int64, track music.Track) (bool, error) {
	isExist, err := t.repo.IsTrackExistInPlaylist(track.Id, idPlaylist)
	if err != nil {
		return false, err
	}
	if isExist {
		return false, errors.New("you already have this track in playlist")
	} else {
		return t.repo.AddTrack(idPlaylist, track)
	}
}

func (t *TrackServiceImpl) ExcludeTrack(idPlaylist int64, idTrack string) error {
	return t.repo.ExcludeTrack(idTrack, idPlaylist)
}

func (t *TrackServiceImpl) GetAllTracks(idPlaylist int64) ([]music.Track, error) {
	return t.repo.SelectAllTracks(idPlaylist)
}

func (t *TrackServiceImpl) GetTrackById(idTrack string) (music.Track, error) {
	return t.repo.SelectTrackById(idTrack)
}

func (t *TrackServiceImpl) GetTrackFromAlbum(idAlbum string) ([]music.Track, error) {
	return t.repo.SelectTracksFromAlbum(idAlbum)
}

func (t *TrackServiceImpl) AddTrackToDB(track music.Track) (bool, error) {
	tx, err := t.repo.StartTransaction()
	if err != nil {
		return false, err
	}
	ok, err := t.repo.InsertTrackToDB(tx, track)
	if err != nil {
		tx.Rollback()
		return false, err
	}
	return ok, tx.Commit()
}

func (t *TrackServiceImpl) GetTracksForRec(userId int64) ([]music.Track, error) {
	return t.repo.SelectTracksByGenre(userId)
}

func (t *TrackServiceImpl) GetTopTracksOfArtist(artistId string) ([]music.Track, error) {
	return t.repo.SelectTopTracksOfArtist(artistId)
}

func (t *TrackServiceImpl) GetMostPopularTracksForUser(userid int64) ([]music.Track, error) {
	return t.repo.SelectMostPopularTracksFroUser(userid)
}
