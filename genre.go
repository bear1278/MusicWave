package music

import (
	"encoding/base64"
	"strconv"
)

type Genre struct {
	Id         int64
	Name       string  `json:"name" binding:"required"`
	Popularity int64   `json:"popularity"`
	Diversity  float64 `json:"diversity"`
}

type GenreJSON struct {
	Id   string `json:"id"`
	Name string `json:"name" binding:"required"`
}

func (g *Genre) SetFromJSON(genre GenreJSON) error {
	decoded, err := base64.StdEncoding.DecodeString(genre.Id)
	if err != nil {
		return err
	}
	g.Id, err = strconv.ParseInt(string(decoded), 10, 64)
	if err != nil {
		return err
	}
	decoded, err = base64.StdEncoding.DecodeString(genre.Name)
	if err != nil {
		return err
	}
	g.Name = string(decoded)
	return err
}
