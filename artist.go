package music

import (
	"encoding/base64"
	"github.com/zmb3/spotify"
	"strconv"
	"time"
)

type Artist struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	Picture    string    `json:"picture"`
	SpotifyURL string    `json:"spotifyURL"`
	Popularity int64     `json:"popularity"`
	AddedDate  time.Time `json:"addedDate"`
	Genres     []Genre   `json:"genres"`
}

type ArtistJSON struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
	SpotifyURL string `json:"spotify_url"`
	Popularity string `json:"popularity"`
	AddedDate  string `json:"added_date"`
}

func (a *Artist) NewArtist(artist spotify.FullArtist) {
	var genres []Genre
	a.Id = artist.ID.String()
	a.Name = artist.Name
	if len(artist.Images) >= 2 {
		a.Picture = artist.Images[1].URL
	} else {
		a.Picture = "/static/images/user.svg"
	}
	a.SpotifyURL = artist.ExternalURLs["spotify"]
	a.Popularity = int64(artist.Popularity)
	for _, g := range artist.Genres {
		genres = append(genres, Genre{
			Id:   0,
			Name: g})
	}
	a.Genres = genres
}

func (a *Artist) SetFromJSON(artist ArtistJSON) error {
	decoded, err := base64.StdEncoding.DecodeString(artist.Id)
	if err != nil {
		return err
	}
	a.Id = string(decoded)
	decoded, err = base64.StdEncoding.DecodeString(artist.Name)
	if err != nil {
		return err
	}
	a.Name = string(decoded)
	decoded, err = base64.StdEncoding.DecodeString(artist.Picture)
	if err != nil {
		return err
	}
	a.Picture = string(decoded)

	decoded, err = base64.StdEncoding.DecodeString(artist.SpotifyURL)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	decoded, err = base64.StdEncoding.DecodeString(artist.Popularity)
	if err != nil {
		return err
	}
	a.Popularity, err = strconv.ParseInt(string(decoded), 10, 64)
	if err != nil {
		return err
	}
	decoded, err = base64.StdEncoding.DecodeString(artist.AddedDate)
	if err != nil {
		return err
	}
	a.AddedDate, err = time.Parse("2006-01-02 15:04:05", string(decoded))
	if err != nil {
		return err
	}
	return err
}
