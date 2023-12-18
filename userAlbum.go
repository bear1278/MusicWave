package music

import (
	"encoding/base64"
	"strconv"
)

type UserAlbum struct {
	IdAlbum string
	IdUser  int64
}

type UserAlbumJSON struct {
	IdArtist string `json:"ID_Album"`
	IdUser   string `json:"ID_User"`
}

func (u *UserAlbum) SetFromJSON(ua UserAlbumJSON) error {
	decoded, err := base64.StdEncoding.DecodeString(ua.IdArtist)
	if err != nil {
		return err
	}
	u.IdAlbum = string(decoded)

	decoded, err = base64.StdEncoding.DecodeString(ua.IdUser)
	if err != nil {
		return err
	}
	u.IdUser, err = strconv.ParseInt(string(decoded), 10, 64)
	if err != nil {
		return err
	}
	return err
}
