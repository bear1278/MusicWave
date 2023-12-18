package music

import (
	"encoding/base64"
	"github.com/zmb3/spotify"
	"strconv"
	"time"
)

type Album struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Cover       string    `json:"cover"`
	ReleaseDate time.Time `json:"releaseDate"`
	AlbumType   string    `json:"albumType"`
	TotalTracks int64     `json:"totalTracks"`
	SpotifyURL  string    `json:"spotifyURL"`
	Popularity  int64     `json:"popularity"`
	Artist      []Artist  `json:"artists"`
	Duration    int64     `json:"duration"`
}

type AlbumJSON struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Cover       string `json:"cover"`
	ReleaseDate string `json:"release_date"`
	AlbumType   string `json:"album_type"`
	TotalTracks string `json:"total_tracks"`
	SpotifyURL  string `json:"spotify_url"`
	Popularity  string `json:"popularity"`
}

func (a *Album) NewAlbum(album spotify.SimpleAlbum) {
	var artist Artist
	a.Id = album.ID.String()
	a.Name = album.Name
	if len(album.Images) >= 2 {
		a.Cover = album.Images[1].URL
	} else {
		a.Cover = "/static/images/music.svg"
	}
	a.ReleaseDate = album.ReleaseDateTime()
	a.AlbumType = album.AlbumType
	a.SpotifyURL = album.ExternalURLs["spotify"]
	for _, artistAPI := range album.Artists {
		artist.Id = artistAPI.ID.String()
		artist.Name = artistAPI.Name
		a.Artist = append(a.Artist, artist)
	}
}

func (a *Album) SetFromJSON(album AlbumJSON) error {
	decoded, err := base64.StdEncoding.DecodeString(album.Id)
	if err != nil {
		return err
	}
	a.Id = string(decoded)
	decoded, err = base64.StdEncoding.DecodeString(album.Name)
	if err != nil {
		return err
	}
	a.Name = string(decoded)
	decoded, err = base64.StdEncoding.DecodeString(album.Cover)
	if err != nil {
		return err
	}
	a.Cover = string(decoded)
	decoded, err = base64.StdEncoding.DecodeString(album.AlbumType)
	if err != nil {
		return err
	}
	a.AlbumType = string(decoded)
	decoded, err = base64.StdEncoding.DecodeString(album.SpotifyURL)
	if err != nil {
		return err
	}
	a.SpotifyURL = string(decoded)
	decoded, err = base64.StdEncoding.DecodeString(album.TotalTracks)
	if err != nil {
		return err
	}
	a.TotalTracks, err = strconv.ParseInt(string(decoded), 10, 64)
	if err != nil {
		return err
	}
	decoded, err = base64.StdEncoding.DecodeString(album.Popularity)
	if err != nil {
		return err
	}
	a.Popularity, err = strconv.ParseInt(string(decoded), 10, 64)
	if err != nil {
		return err
	}
	decoded, err = base64.StdEncoding.DecodeString(album.ReleaseDate)
	if err != nil {
		return err
	}
	a.ReleaseDate, err = time.Parse("2006-01-02 15:04:05", string(decoded))
	if err != nil {
		return err
	}
	return err
}
