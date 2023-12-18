package music

import (
	"encoding/base64"
	"strconv"
)

type UserPlaylist struct {
	IdPlaylist int64
	IdUser     int64
}

type UserPlaylistJSON struct {
	IdPlaylist string `json:"ID_Playlist"`
	IdUser     string `json:"ID_User"`
}

func (u *UserPlaylist) SetFromJSON(up UserPlaylistJSON) error {
	decoded, err := base64.StdEncoding.DecodeString(up.IdPlaylist)
	if err != nil {
		return err
	}
	u.IdPlaylist, err = strconv.ParseInt(string(decoded), 10, 64)
	if err != nil {
		return err
	}

	decoded, err = base64.StdEncoding.DecodeString(up.IdUser)
	if err != nil {
		return err
	}
	u.IdUser, err = strconv.ParseInt(string(decoded), 10, 64)
	if err != nil {
		return err
	}
	return err
}
