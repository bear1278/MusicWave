package repository

import (
	"database/sql"
	"fmt"
	music "github.com/bear1278/MusicWave"
	"time"
)

type PlaylistMysql struct {
	db *sql.DB
}

func NewPlaylistMysql(db *sql.DB) *PlaylistMysql {
	return &PlaylistMysql{db: db}
}

func (p *PlaylistMysql) CreatePlaylist(playlist music.Playlist) (int64, error) {
	transaction, err := p.db.Begin()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (name,duration,cover,author,release_date) values (?,?,?,?,?)", playlistTable)
	dbResult, err := transaction.Exec(query, playlist.Name, playlist.Duration, playlist.Cover, playlist.Author.Id, playlist.ReleaseDate)
	if err != nil {
		err = transaction.Rollback()
		return 0, err
	}
	idPlaylist, err := dbResult.LastInsertId()
	if err != nil {
		return 0, err
	}
	query = fmt.Sprintf("INSERT INTO %s values (?,?)", userPlaylistTable)
	_, err = transaction.Exec(query, playlist.Author.Id, idPlaylist)
	if err != nil {
		err = transaction.Rollback()
		return 0, err
	}
	if err := transaction.Commit(); err != nil {
		err = transaction.Rollback()
		return 0, err
	}
	return idPlaylist, nil
}

func (p *PlaylistMysql) DeletePlaylist(playlist int64) error {
	transaction, err := p.db.Begin()
	if err != nil {
		return err
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE ID_Playlist=?", userPlaylistTable)
	_, err = transaction.Exec(query, playlist)
	if err != nil {
		err = transaction.Rollback()
		return err
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE id=?", playlistTable)
	_, err = transaction.Exec(query, playlist)
	if err != nil {
		err = transaction.Rollback()
		return err
	}
	if err := transaction.Commit(); err != nil {
		err = transaction.Rollback()
		return err
	}
	return nil
}

func (p *PlaylistMysql) SelectAllPlaylists(idUser int64) ([]music.Playlist, error) {
	var playlists []music.Playlist
	var playlist music.Playlist
	var date string
	query := fmt.Sprintf("SELECT P.id, P.name, duration, cover, U.id, U.Username, release_date FROM %s as UP INNER JOIN %s as P ON P.id=UP.ID_Playlist "+
		"INNER JOIN Users as U ON U.id=UP.ID_User WHERE UP.ID_User=?", userPlaylistTable, playlistTable)
	rows, err := p.db.Query(query, idUser)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&playlist.Id, &playlist.Name, &playlist.Duration, &playlist.Cover, &playlist.Author.Id, &playlist.Author.UserName, &date)
		if err != nil {
			return nil, err
		}
		playlist.ReleaseDate, err = time.Parse(time.DateTime, date)
		if err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}
	return playlists, nil
}

func (p *PlaylistMysql) SelectById(playlistId int64) (music.Playlist, error) {
	var playlist music.Playlist
	var date string
	query := fmt.Sprintf("SELECT P.id, P.name, duration, cover, U.id, U.Username, release_date  FROM %s as P INNER JOIN %s as U on U.id=P.Author WHERE P.id=?", playlistTable, usersTable)
	rows, err := p.db.Query(query, playlistId)
	if err != nil {
		return playlist, err
	}
	for rows.Next() {
		err = rows.Scan(&playlist.Id, &playlist.Name, &playlist.Duration, &playlist.Cover, &playlist.Author.Id, &playlist.Author.UserName, &date)
		if err != nil {
			return playlist, err
		}
		playlist.ReleaseDate, err = time.Parse(time.DateTime, date)
		if err != nil {
			return playlist, err
		}
	}
	return playlist, nil
}

func (p *PlaylistMysql) UpdatePlaylist(playlist music.Playlist) error {
	query := fmt.Sprintf("UPDATE %s SET name=?, cover=? where id=?", playlistTable)
	_, err := p.db.Exec(query, playlist.Name, playlist.Cover, playlist.Id)
	return err
}

func (p *PlaylistMysql) ExcludePlaylist(idPlaylist, idUser int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE ID_playlist=? AND ID_user=?", userPlaylistTable)
	_, err := p.db.Exec(query, idPlaylist, idUser)
	return err
}

func (p *PlaylistMysql) AddPlaylist(idPlaylist, idUser int64) error {
	query := fmt.Sprintf("INSERT INTO %s values (?,?)", userPlaylistTable)
	_, err := p.db.Exec(query, idUser, idPlaylist)
	return err
}
