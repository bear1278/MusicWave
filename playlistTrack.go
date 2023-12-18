package music

import (
	"encoding/base64"
	"strconv"
)

type PlaylistTrack struct {
	IdTrack    string
	IdPlaylist int64
}

type PlaylistTrackJSON struct {
	IdTrack    string `json:"ID_Track"`
	IdPlaylist string `json:"ID_Playlist"`
}

func (p *PlaylistTrack) SetFromJSON(pt PlaylistTrackJSON) error {
	decoded, err := base64.StdEncoding.DecodeString(pt.IdTrack)
	if err != nil {
		return err
	}
	p.IdTrack = string(decoded)

	decoded, err = base64.StdEncoding.DecodeString(pt.IdPlaylist)
	if err != nil {
		return err
	}
	p.IdPlaylist, err = strconv.ParseInt(string(decoded), 10, 64)
	if err != nil {
		return err
	}
	return err
}
