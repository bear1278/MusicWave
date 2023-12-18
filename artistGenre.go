package music

import (
	"encoding/base64"
	"strconv"
)

type ArtistGenre struct {
	IdArtist string
	IdGenre  int64
}

type ArtistGenreJSON struct {
	IdArtist string `json:"ID_Artist"`
	IdGenre  string `json:"ID_Genre"`
}

func (a *ArtistGenre) SetFromJSON(ag ArtistGenreJSON) error {
	decoded, err := base64.StdEncoding.DecodeString(ag.IdArtist)
	if err != nil {
		return err
	}
	a.IdArtist = string(decoded)

	decoded, err = base64.StdEncoding.DecodeString(ag.IdGenre)
	if err != nil {
		return err
	}
	a.IdGenre, err = strconv.ParseInt(string(decoded), 10, 64)
	if err != nil {
		return err
	}
	return err
}
