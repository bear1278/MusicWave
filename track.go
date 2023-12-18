package music

import (
	"encoding/base64"
	"github.com/zmb3/spotify"
	"strconv"
)

type Track struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Duration   int64  `json:"duration"`
	Cover      string `json:"cover"`
	Album      Album  `json:"album"`
	SpotifyURL string `json:"spotifyURL"`
	Popularity int64  `json:"popularity"`
}

type TrackJSON struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Duration   string `json:"duration"`
	Cover      string `json:"cover"`
	Album      string `json:"ID_Album"`
	SpotifyURL string `json:"spotify_url"`
	Popularity string `json:"popularity"`
}

func (t *Track) NewTrack(track spotify.FullTrack) {
	var artist Artist
	t.Id = track.ID.String()
	t.Name = track.Name
	t.Duration = int64(track.Duration)
	if len(track.Album.Images) >= 2 {
		t.Cover = track.Album.Images[1].URL
	} else {
		t.Cover = "/static/images/music.svg"
	}
	t.SpotifyURL = track.ExternalURLs["spotify"]
	t.Popularity = int64(track.Popularity)
	t.Album.Id = track.Album.ID.String()
	t.Album.Name = track.Name
	for _, artistAPI := range track.Artists {
		artist.Id = artistAPI.ID.String()
		artist.Name = artistAPI.Name
		t.Album.Artist = append(t.Album.Artist, artist)
	}
}

func (t *Track) SetFromJSON(track TrackJSON) error {
	decoded, err := base64.StdEncoding.DecodeString(track.Id)
	if err != nil {
		return err
	}
	t.Id = string(decoded)
	decoded, err = base64.StdEncoding.DecodeString(track.Name)
	if err != nil {
		return err
	}
	t.Name = string(decoded)
	decoded, err = base64.StdEncoding.DecodeString(track.Cover)
	if err != nil {
		return err
	}
	t.Cover = string(decoded)
	decoded, err = base64.StdEncoding.DecodeString(track.Album)
	if err != nil {
		return err
	}
	t.Album.Id = string(decoded)
	decoded, err = base64.StdEncoding.DecodeString(track.SpotifyURL)
	if err != nil {
		return err
	}
	t.SpotifyURL = string(decoded)
	decoded, err = base64.StdEncoding.DecodeString(track.Duration)
	if err != nil {
		return err
	}
	t.Duration, err = strconv.ParseInt(string(decoded), 10, 64)
	if err != nil {
		return err
	}
	decoded, err = base64.StdEncoding.DecodeString(track.Popularity)
	if err != nil {
		return err
	}
	t.Popularity, err = strconv.ParseInt(string(decoded), 10, 64)
	if err != nil {
		return err
	}
	return err
}
