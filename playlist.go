package music

import "time"

type Playlist struct {
	Id          int64  `json:"id"`
	Name        string `json:"name" binding:"required"`
	Duration    int64
	Cover       string `json:"cover" binding:"required"`
	Author      User
	ReleaseDate time.Time
}
