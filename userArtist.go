package music

import (
	"encoding/base64"
	"strconv"
)

type UserArtist struct {
	IdArtist string
	IdUser   int64
}

type UserArtistJSON struct {
	IdArtist string `json:"ID_Artist"`
	IdUser   string `json:"ID_User"`
}

func (u *UserArtist) SetFromJSON(ua UserArtistJSON) error {
	decoded, err := base64.StdEncoding.DecodeString(ua.IdArtist)
	if err != nil {
		return err
	}
	u.IdArtist = string(decoded)

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
