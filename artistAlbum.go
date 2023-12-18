package music

import (
	"encoding/base64"
)

type ArtistAlbum struct {
	IdArtist string
	IdAlbum  string
}

type ArtistAlbumJSON struct {
	IdArtist string `json:"ID_Artist"`
	IdAlbum  string `json:"ID_Album"`
}

func (a *ArtistAlbum) SetFromJSON(aa ArtistAlbumJSON) error {
	decoded, err := base64.StdEncoding.DecodeString(aa.IdArtist)
	if err != nil {
		return err
	}
	a.IdArtist = string(decoded)

	decoded, err = base64.StdEncoding.DecodeString(aa.IdAlbum)
	if err != nil {
		return err
	}
	a.IdAlbum = string(decoded)
	return err
}
