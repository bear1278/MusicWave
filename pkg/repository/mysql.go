package repository

import (
	"database/sql"
	"fmt"
	"github.com/bear1278/MusicWave/configs"
	_ "github.com/go-sql-driver/mysql"
)

const (
	usersTable         = "users"
	genresTable        = "genres"
	userGenreTable     = "user_genre"
	playlistTable      = "playlists"
	userPlaylistTable  = "user_playlist"
	trackTable         = "tracks"
	playlistTrackTable = "playlist_track"
	albumTable         = "Albums"
	artistTable        = "Artists"
	artistAlbumTable   = "Artist_album"
	artistGenreTable   = "Artist_Genre"
	userAlbumTable     = "User_Album"
	userArtistTable    = "User_Artist"
	adminHistory       = "adminhistory"
)

func MySqlDB(cfg configs.Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.DB.Driver, fmt.Sprintf("%s:%s@/%s", cfg.DB.User, cfg.DB.Password, cfg.DB.DbName))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
