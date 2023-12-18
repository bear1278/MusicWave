package music

import (
	"encoding/base64"
	"strconv"
	"time"
)

type Playlist struct {
	Id           int64  `json:"id"`
	Name         string `json:"name" binding:"required"`
	Duration     int64  `json:"duration"`
	Cover        string `json:"cover" binding:"required"`
	Author       User
	ReleaseDate  time.Time `json:"releaseDate"`
	Type         string    `json:"type"`
	ModifiedDate time.Time `json:"modifiedDate"`
}

type PlaylistJSON struct {
	Id           string `json:"id"`
	Name         string `json:"name" binding:"required"`
	Duration     string `json:"duration"`
	Cover        string `json:"cover" binding:"required"`
	Author       string `json:"author"`
	ReleaseDate  string `json:"release_date"`
	Type         string `json:"type"`
	ModifiedDate string `json:"modified_date"`
}

func (p *Playlist) SetFromJSON(playlist PlaylistJSON) error {
	decoded, err := base64.StdEncoding.DecodeString(playlist.Id)
	if err != nil {
		return err
	}
	p.Id, err = strconv.ParseInt(string(decoded), 10, 64)
	if err != nil {
		return err
	}
	decoded, err = base64.StdEncoding.DecodeString(playlist.Name)
	if err != nil {
		return err
	}
	p.Name = string(decoded)
	decoded, err = base64.StdEncoding.DecodeString(playlist.Cover)
	if err != nil {
		return err
	}
	p.Cover = string(decoded)
	decoded, err = base64.StdEncoding.DecodeString(playlist.Type)
	if err != nil {
		return err
	}
	p.Type = string(decoded)

	decoded, err = base64.StdEncoding.DecodeString(playlist.Duration)
	if err != nil {
		return err
	}
	p.Duration, err = strconv.ParseInt(string(decoded), 10, 64)
	if err != nil {
		return err
	}
	decoded, err = base64.StdEncoding.DecodeString(playlist.Author)
	if err != nil {
		return err
	}
	p.Author.Id, err = strconv.ParseInt(string(decoded), 10, 64)
	if err != nil {
		return err
	}
	decoded, err = base64.StdEncoding.DecodeString(playlist.ReleaseDate)
	if err != nil {
		return err
	}
	p.ReleaseDate, err = time.Parse("2006-01-02 15:04:05", string(decoded))
	if err != nil {
		return err
	}
	decoded, err = base64.StdEncoding.DecodeString(playlist.ModifiedDate)
	if err != nil {
		return err
	}
	p.ModifiedDate, err = time.Parse("2006-01-02 15:04:05", string(decoded))
	if err != nil {
		return err
	}
	return err
}
