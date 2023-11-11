package repository

import (
	"database/sql"
)

type TrackMysql struct {
	db *sql.DB
}

func NewTrackMysql(db *sql.DB) *TrackMysql {
	return &TrackMysql{db: db}
}

//func (t *TrackMysql)AddTrack(idPlaylist int64, track music.Track)(int64,error){
//
//}
