package repository

import (
	"database/sql"
	"errors"
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

func (p *PlaylistMysql) IsExistInUserLibrary(playlistId, userId int64) (bool, error) {
	fav, err := p.GetUserFavourites(userId)
	if err != nil {
		return false, err
	}
	if fav.Id == playlistId {
		return true, err
	}
	query := fmt.Sprintf("SELECT * FROM %s WHERE ID_Playlist=? and ID_User=?", userPlaylistTable)
	rows, err := p.db.Query(query, playlistId, userId)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, err
	} else {
		return false, err
	}
}

func (p *PlaylistMysql) CreatePlaylist(playlist music.Playlist) (int64, error) {
	transaction, err := p.db.Begin()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (name,duration,cover,author,release_date,type,modified_date) values (?,?,?,?,?,?,?)", playlistTable)
	dbResult, err := transaction.Exec(query, playlist.Name, playlist.Duration, playlist.Cover,
		playlist.Author.Id, playlist.ReleaseDate, playlist.Type, time.Now())
	if err != nil {
		transaction.Rollback()
		return 0, err
	}
	idPlaylist, err := dbResult.LastInsertId()
	if err != nil {
		return 0, err
	}
	query = fmt.Sprintf("INSERT INTO %s values (?,?)", userPlaylistTable)
	_, err = transaction.Exec(query, playlist.Author.Id, idPlaylist)
	if err != nil {
		transaction.Rollback()
		return 0, err
	}
	if err := transaction.Commit(); err != nil {
		transaction.Rollback()
		return 0, err
	}
	return idPlaylist, nil
}

func (p *PlaylistMysql) DeletePlaylist(playlistID, userID int64) error {
	transaction, err := p.db.Begin()
	if err != nil {
		return err
	}
	isExist, err := p.CheckIfExist(playlistID)
	if err != nil {
		return err
	}
	if isExist {
		query := fmt.Sprintf("SELECT * FROM %s WHERE id=? AND author=?", playlistTable)
		rows, err := p.db.Query(query, playlistID, userID)
		if err != nil {
			return err
		}
		if !rows.Next() {
			return errors.New("you are not author of this playlist")
		}
		query = fmt.Sprintf("DELETE FROM %s WHERE ID_Playlist=?", userPlaylistTable)
		_, err = transaction.Exec(query, playlistID)
		if err != nil {
			transaction.Rollback()
			return err
		}
		query = fmt.Sprintf("DELETE FROM %s WHERE id=?", playlistTable)
		_, err = transaction.Exec(query, playlistID)
		if err != nil {
			transaction.Rollback()
			return err
		}
		return transaction.Commit()
	} else {
		return errors.New("playlist doesn't exist")
	}
}

func (p *PlaylistMysql) SelectAllPlaylists(idUser int64) ([]music.Playlist, error) {
	var playlists []music.Playlist
	var playlist music.Playlist
	var date string
	query := fmt.Sprintf("select P.id, P.name, P.duration, P.cover, P.uid, P.Username, P.release_date, P.type from "+
		"(SELECT P.id, P.name, duration, cover, U.id as uid, U.Username, release_date, P.type FROM %s as P "+
		"INNER JOIN %s as U ON U.id=P.author) as P "+
		"INNER JOIN %s as UP on UP.ID_Playlist=P.id "+
		"WHERE UP.ID_User=?", playlistTable, usersTable, userPlaylistTable)
	rows, err := p.db.Query(query, idUser)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&playlist.Id, &playlist.Name, &playlist.Duration, &playlist.Cover, &playlist.Author.Id, &playlist.Author.UserName, &date, &playlist.Type)
		if err != nil {
			return nil, err
		}
		playlist.ReleaseDate, err = time.Parse(time.DateTime, date)
		if err != nil {
			return nil, err
		}
		if playlist.Cover == "" || playlist.Cover == " " || playlist.Cover == "nil" {
			playlist.Cover = "/static/images/music.svg"
		}
		playlists = append(playlists, playlist)
	}
	return playlists, nil
}

func (p *PlaylistMysql) SelectById(playlistId int64) (music.Playlist, error) {
	var playlist music.Playlist
	var date string
	isExist, err := p.CheckIfExist(playlistId)
	if err != nil {
		return music.Playlist{}, err
	}
	if isExist {
		query := fmt.Sprintf("SELECT P.id, P.name, duration, cover, U.id, U.Username, release_date, type  FROM %s as P "+
			"INNER JOIN %s as U on U.id=P.Author WHERE P.id=?", playlistTable, usersTable)
		rows, err := p.db.Query(query, playlistId)
		if err != nil {
			return playlist, err
		}
		for rows.Next() {
			err = rows.Scan(&playlist.Id, &playlist.Name, &playlist.Duration, &playlist.Cover, &playlist.Author.Id, &playlist.Author.UserName, &date, &playlist.Type)
			if err != nil {
				return playlist, err
			}
			playlist.ReleaseDate, err = time.Parse(time.DateTime, date)
			if playlist.Cover == "" || playlist.Cover == " " || playlist.Cover == "nil" {
				playlist.Cover = "/static/images/music.svg"
			}
			if err != nil {
				return playlist, err
			}
		}
		return playlist, nil
	} else {
		return music.Playlist{}, errors.New("playlist doesn't exist")
	}
}

func (p *PlaylistMysql) UpdatePlaylist(playlist music.Playlist) error {
	isExist, err := p.CheckIfExist(playlist.Id)
	if err != nil {
		return err
	}
	if isExist {
		query := fmt.Sprintf("SELECT * FROM %s WHERE id=? AND author=?", playlistTable)
		rows, err := p.db.Query(query, playlist.Id, playlist.Author.Id)
		if err != nil {
			return err
		}
		if !rows.Next() {
			return errors.New("you are not author of this playlist")
		}
		query = fmt.Sprintf("UPDATE %s SET name=?, cover=?, type=?,modified_date=? where id=?", playlistTable)
		_, err = p.db.Exec(query, playlist.Name, playlist.Cover, playlist.Type, time.Now(), playlist.Id)
		return err
	} else {
		return errors.New("playlist doesn't exist")
	}
}

func (p *PlaylistMysql) ExcludePlaylist(idPlaylist, idUser int64) error {
	isExist, err := p.CheckIfExist(idPlaylist)
	if err != nil {
		return err
	}
	if isExist {
		query := fmt.Sprintf("DELETE FROM %s WHERE ID_playlist=? AND ID_user=?", userPlaylistTable)
		_, err = p.db.Exec(query, idPlaylist, idUser)
		return err
	} else {
		return errors.New("playlist doesn't exist")
	}

}

func (p *PlaylistMysql) CheckIfExist(idPlaylist int64) (bool, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=?", playlistTable)
	rows, err := p.db.Query(query, idPlaylist)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, err
	} else {
		return false, err
	}
}

func (p *PlaylistMysql) AddPlaylist(idPlaylist, idUser int64) error {
	isExist, err := p.CheckIfExist(idPlaylist)
	if err != nil {
		return err
	}
	if isExist {
		query := fmt.Sprintf("INSERT INTO %s values (?,?)", userPlaylistTable)
		_, err = p.db.Exec(query, idUser, idPlaylist)
		return err
	} else {
		return errors.New("playlist doesn't exist")
	}
}

func (p *PlaylistMysql) SelectPlaylistForSearch(search string) ([]music.Playlist, error) {
	var playlist music.Playlist
	var playlists []music.Playlist
	var date string
	query := fmt.Sprintf("SELECT P.id, P.name, duration, cover, U.id, U.Username, release_date, P.type FROM %s as P "+
		"INNER JOIN %s as U ON U.id=P.Author "+
		"WHERE P.name!='favorites' and P.type!='private' and P.name REGEXP ?", playlistTable, usersTable)
	rows, err := p.db.Query(query, ".*"+search+".*")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&playlist.Id, &playlist.Name, &playlist.Duration, &playlist.Cover, &playlist.Author.Id, &playlist.Author.UserName, &date, &playlist.Type)
		if err != nil {
			return nil, err
		}
		playlist.ReleaseDate, err = time.Parse(time.DateTime, date)
		if err != nil {
			return nil, err
		}
		if playlist.Cover == "" || playlist.Cover == " " || playlist.Cover == "nil" {
			playlist.Cover = "/static/images/music.svg"
		}
		playlists = append(playlists, playlist)
	}
	return playlists, nil
}

func (p *PlaylistMysql) SelectAllPlaylistsForAdmin() ([]music.Playlist, error) {
	var playlist music.Playlist
	var playlists []music.Playlist
	var releaseDate, modDate string
	query := fmt.Sprintf("SELECT P.id, P.name, duration, cover, U.id, U.Username, release_date, P.type,P.modified_date FROM %s as P "+
		"INNER JOIN %s as U on P.Author=U.id "+
		"WHERE P.name!='favorites'", playlistTable, usersTable)
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&playlist.Id, &playlist.Name, &playlist.Duration,
			&playlist.Cover, &playlist.Author.Id, &playlist.Author.UserName, &releaseDate, &playlist.Type, &modDate)
		if err != nil {
			return nil, err
		}
		playlist.ReleaseDate, err = time.Parse(time.DateTime, releaseDate)
		if err != nil {
			return nil, err
		}
		playlist.ModifiedDate, err = time.Parse(time.DateTime, modDate)
		if playlist.Cover == "" || playlist.Cover == " " || playlist.Cover == "nil" {
			playlist.Cover = "/static/images/music.svg"
		}
		playlists = append(playlists, playlist)
		playlist = music.Playlist{}
	}
	return playlists, err
}

func (p *PlaylistMysql) DeletePlaylistByAdmin(playlistId int64, reason string) error {
	var name string
	transaction, err := p.db.Begin()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT name FROM %s WHERE id=?", playlistTable)
	rows, err := transaction.Query(query, playlistId)
	if err != nil {
		transaction.Rollback()
		return err
	}
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			transaction.Rollback()
			return err
		}
	}
	query = fmt.Sprintf("INSERT INTO %s (name,type,reason,deleted_date) values (?,?,?,?)", adminHistory)
	_, err = transaction.Exec(query, name, "playlist", reason, time.Now())
	if err != nil {
		transaction.Rollback()
		return err
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE ID_Playlist=?", userPlaylistTable)
	_, err = transaction.Exec(query, playlistId)
	if err != nil {
		transaction.Rollback()
		return err
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE id=?", playlistTable)
	_, err = transaction.Exec(query, playlistId)
	if err != nil {
		transaction.Rollback()
		return err
	}

	return transaction.Commit()
}

func (p *PlaylistMysql) GetUserFavourites(userID int64) (music.Playlist, error) {
	var favorites music.Playlist
	var createDate string
	query := fmt.Sprintf("SELECT P.id, P.name, duration, cover, U.id, U.Username, release_date, type  FROM %s as P "+
		"INNER JOIN %s as U on U.id=P.Author WHERE U.id=? and P.name=?", playlistTable, usersTable)
	rows, err := p.db.Query(query, userID, "favorites")
	if err != nil {
		return music.Playlist{}, err
	}
	for rows.Next() {
		err = rows.Scan(&favorites.Id, &favorites.Name, &favorites.Duration, &favorites.Cover,
			&favorites.Author.Id, &favorites.Author.UserName, &createDate, &favorites.Type)
		if err != nil {
			return music.Playlist{}, err
		}
		favorites.ReleaseDate, err = time.Parse(time.DateTime, createDate)
		if err != nil {
			return music.Playlist{}, err
		}
		if favorites.Cover == "" || favorites.Cover == " " || favorites.Cover == "nil" {
			favorites.Cover = "/static/images/music.svg"
		}

	}
	return favorites, err
}

func (p *PlaylistMysql) SelectPlaylistsRec(userId int64) ([]music.Playlist, error) {
	var playlist music.Playlist
	var playlists []music.Playlist
	var releaseDate, modDate string
	query := fmt.Sprintf("SELECT P.id, P.name, duration, cover, U.id, U.Username, release_date, P.type,P.modified_date FROM %s as P "+
		"INNER JOIN %s as U on P.Author=U.id "+
		"WHERE P.type!='private' and P.name!='favorites' and P.Author!=? and p.id not in "+
		"(select P1.id from %s as P1 inner join %s as UP1 on P1.id=UP1.ID_Playlist where UP1.ID_User=?) "+
		"ORDER BY P.release_date desc "+
		"LIMIT 10", playlistTable, usersTable, playlistTable, userPlaylistTable)
	rows, err := p.db.Query(query, userId, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&playlist.Id, &playlist.Name, &playlist.Duration,
			&playlist.Cover, &playlist.Author.Id, &playlist.Author.UserName, &releaseDate, &playlist.Type, &modDate)
		if err != nil {
			return nil, err
		}
		if playlist.Cover == "" || playlist.Cover == " " || playlist.Cover == "nil" {
			playlist.Cover = "/static/images/music.svg"
		}
		playlist.ReleaseDate, err = time.Parse(time.DateTime, releaseDate)
		if err != nil {
			return nil, err
		}
		playlist.ModifiedDate, err = time.Parse(time.DateTime, modDate)
		playlists = append(playlists, playlist)
		playlist = music.Playlist{}
	}
	return playlists, err
}

func (p *PlaylistMysql) SelectPlaylistDuration(playlistId int64) (int64, error) {
	var duration sql.NullInt64
	query := fmt.Sprintf("Select sum(T.duration) as duration from %s as PT "+
		"INNER JOIN %s AS T on T.id=PT.ID_Track "+
		"WHERE PT.ID_Playlist=?", playlistTrackTable, trackTable)
	rows, err := p.db.Query(query, playlistId)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		err = rows.Scan(&duration)
		if err != nil {
			return 0, err
		}
	}
	if duration.Valid {
		return duration.Int64, err
	} else {
		return 0, err
	}
}

func (p *PlaylistMysql) SelectIsAddedToSpotify(playlistId int64) (bool, error) {
	var isAdded bool
	query := fmt.Sprintf("Select spotify_added FROM %s where id=?", playlistTable)
	rows, err := p.db.Query(query, playlistId)
	if err != nil {
		return false, err
	}
	for rows.Next() {
		err = rows.Scan(&isAdded)
		if err != nil {
			return false, err
		}
	}
	return isAdded, err
}

func (p *PlaylistMysql) UpdateIsAdded(playlistId int64) error {
	query := fmt.Sprintf("UPDATE %s SET spotify_added=1 where id=?", playlistTable)
	_, err := p.db.Exec(query, playlistId)
	return err
}
