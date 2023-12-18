package music

import (
	"encoding/base64"
	"strconv"
)

type UserGenre struct {
	IdGenre int64
	IdUser  int64
}

type UserGenreJSON struct {
	IdGenre string `json:"ID_Genre"`
	IdUser  string `json:"ID_User"`
}

func (u *UserGenre) SetFromJSON(ug UserGenreJSON) error {
	decoded, err := base64.StdEncoding.DecodeString(ug.IdGenre)
	if err != nil {
		return err
	}
	u.IdGenre, err = strconv.ParseInt(string(decoded), 10, 64)
	if err != nil {
		return err
	}

	decoded, err = base64.StdEncoding.DecodeString(ug.IdUser)
	if err != nil {
		return err
	}
	u.IdUser, err = strconv.ParseInt(string(decoded), 10, 64)
	if err != nil {
		return err
	}
	return err
}
