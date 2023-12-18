package service

import (
	"errors"
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

func (p *PlaylistServiceImpl) DeletePlaylist(playlistID, userID int64) error {
	return p.repo.DeletePlaylist(playlistID, userID)

}

func (p *PlaylistServiceImpl) GetAllPlaylist(idUser int64) ([]music.Playlist, []music.Playlist, error) {
	var myPlaylists, addedPlaylists []music.Playlist
	playlists, err := p.repo.SelectAllPlaylists(idUser)
	if err != nil {
		return nil, nil, err
	}
	for k, _ := range playlists {
		playlists[k].Duration, err = p.repo.SelectPlaylistDuration(playlists[k].Id)
		if err != nil {
			return nil, nil, err
		}
	}
	for _, playlist := range playlists {
		if playlist.Author.Id == idUser {
			myPlaylists = append(myPlaylists, playlist)
		} else {
			addedPlaylists = append(addedPlaylists, playlist)
		}
	}
	return myPlaylists, addedPlaylists, err
}

func (p *PlaylistServiceImpl) GetById(playlistId int64) (music.Playlist, error) {
	playlist, err := p.repo.SelectById(playlistId)
	if err != nil {
		return music.Playlist{}, err
	}
	playlist.Duration, err = p.repo.SelectPlaylistDuration(playlist.Id)
	if err != nil {
		return music.Playlist{}, err
	}
	return playlist, err
}

func (p *PlaylistServiceImpl) UpdatePlaylist(playlist music.Playlist) error {
	return p.repo.UpdatePlaylist(playlist)
}

func (p *PlaylistServiceImpl) ExcludePlaylist(idPlaylist, idUser int64) error {
	return p.repo.ExcludePlaylist(idPlaylist, idUser)
}

func (p *PlaylistServiceImpl) AddPlaylist(idPlaylist, idUser int64) error {
	isExist, err := p.repo.IsExistInUserLibrary(idPlaylist, idUser)
	if err != nil {
		return err
	}
	if isExist {
		return errors.New("you already have this playlist in your library")
	} else {
		return p.repo.AddPlaylist(idPlaylist, idUser)
	}
}

func (p *PlaylistServiceImpl) SearchPlaylist(search string) ([]music.Playlist, error) {
	playlists, err := p.repo.SelectPlaylistForSearch(search)
	if err != nil {
		return nil, err
	}
	for k, _ := range playlists {
		playlists[k].Duration, err = p.repo.SelectPlaylistDuration(playlists[k].Id)
		if err != nil {
			return nil, err
		}
	}
	return playlists, err
}

func (p *PlaylistServiceImpl) GetAllPlaylistForAdmin() ([]music.Playlist, error) {
	playlists, err := p.repo.SelectAllPlaylistsForAdmin()
	if err != nil {
		return nil, err
	}
	for k, _ := range playlists {
		playlists[k].Duration, err = p.repo.SelectPlaylistDuration(playlists[k].Id)
		if err != nil {
			return nil, err
		}
	}
	return playlists, err
}

func (p *PlaylistServiceImpl) DeletePlaylistByAdmin(playlistID int64, reason string) error {
	return p.repo.DeletePlaylistByAdmin(playlistID, reason)
}

func (p *PlaylistServiceImpl) GetUserFavorites(userID int64) (music.Playlist, error) {
	favorites, err := p.repo.GetUserFavourites(userID)
	if err != nil {
		return music.Playlist{}, err
	}
	if favorites.Duration == 0 {
		favorites.Duration, err = p.repo.SelectPlaylistDuration(favorites.Id)
		if err != nil {
			return music.Playlist{}, err
		}
	}
	return favorites, err
}

func (p *PlaylistServiceImpl) GerPlaylistsRec(userId int64) ([]music.Playlist, error) {
	playlists, err := p.repo.SelectPlaylistsRec(userId)
	if err != nil {
		return nil, err
	}
	for k, _ := range playlists {
		playlists[k].Duration, err = p.repo.SelectPlaylistDuration(playlists[k].Id)
		if err != nil {
			return nil, err
		}
	}
	return playlists, err
}

func (p *PlaylistServiceImpl) IsAddedToSpotify(playlistId int64) (bool, error) {
	return p.repo.SelectIsAddedToSpotify(playlistId)
}

func (p *PlaylistServiceImpl) UpdateIsAdded(playlistId int64) error {
	return p.repo.UpdateIsAdded(playlistId)
}

func (p *PlaylistServiceImpl) GetDurationOfAllPlaylists(userId int64) (int64, error) {
	var duration int64 = 0
	playlists, err := p.repo.SelectAllPlaylists(userId)
	if err != nil {
		return 0, err
	}
	for k, _ := range playlists {
		playlists[k].Duration, err = p.repo.SelectPlaylistDuration(playlists[k].Id)
		if err != nil {
			return 0, err
		}
		duration += playlists[k].Duration
	}
	fav, err := p.repo.GetUserFavourites(userId)
	if err != nil {
		return 0, err
	}
	fav.Duration, err = p.repo.SelectPlaylistDuration(fav.Id)
	if err != nil {
		return 0, err
	}
	duration += fav.Duration
	return duration, err
}
