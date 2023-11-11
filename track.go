package music

import "time"

type Track struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Duration    int64     `json:"duration"`
	Cover       string    `json:"cover"`
	ReleaseDate time.Time `json:"release_Date"`
}
