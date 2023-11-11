package service

import (
	music "github.com/bear1278/MusicWave"
	"github.com/bear1278/MusicWave/pkg/repository"
	"time"
)

type PlaylistServiceImpl struct {
	repo repository.PlaylistRepo
}

func NewPlaylistServiceImpl(repo repository.PlaylistRepo) *PlaylistServiceImpl {
	return &PlaylistServiceImpl{repo: repo}
}

func (p *PlaylistServiceImpl) NewPlaylist(playlist music.Playlist, userId int64) (int64, error) {
	playlist.Author.Id = userId
	playlist.ReleaseDate = time.Now()
	playlist.Duration = 0
	idPlaylist, err := p.repo.CreatePlaylist(playlist)
	if err != nil {
		return 0, err
	}
	return idPlaylist, nil
}

func (p *PlaylistServiceImpl) DeletePlaylist(playlist int64) error {
	return p.repo.DeletePlaylist(playlist)

}

func (p *PlaylistServiceImpl) GetAllPlaylist(idUser int64) ([]music.Playlist, error) {
	return p.repo.SelectAllPlaylists(idUser)
}

func (p *PlaylistServiceImpl) GetById(playlist int64) (music.Playlist, error) {
	return p.repo.SelectById(playlist)
}

func (p *PlaylistServiceImpl) UpdatePlaylist(playlist music.Playlist) error {
	return p.repo.UpdatePlaylist(playlist)
}

func (p *PlaylistServiceImpl) ExcludePlaylist(idPlaylist, idUser int64) error {
	return p.repo.ExcludePlaylist(idPlaylist, idUser)
}

func (p *PlaylistServiceImpl) AddPlaylist(idPlaylist, idUser int64) error {
	return p.repo.AddPlaylist(idPlaylist, idUser)
}
