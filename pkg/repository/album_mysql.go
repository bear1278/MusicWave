package repository

import (
	"database/sql"
	"fmt"
	music "github.com/bear1278/MusicWave"
	"time"
)

type AlbumMysql struct {
	db *sql.DB
}

func NewAlbumMysql(db *sql.DB) *AlbumMysql {
	return &AlbumMysql{db: db}
}

func (a *AlbumMysql) IsAlbumExistInUserLibrary(albumId string, userId int64) (bool, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE ID_Album=? and ID_User=?", userAlbumTable)
	rows, err := a.db.Query(query, albumId, userId)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, err
	} else {
		return false, err
	}
}

func (a *AlbumMysql) GetAlbumById(idAlbum string) (music.Album, error) {
	var album music.Album
	var date, artistId, artistName string
	var artist music.Artist
	query := fmt.Sprintf("SELECT A.id,A.name,A.cover,A.release_date,ART.id,ART.name,A.album_type,A.total_tracks,A.spotify_URL,A.popularity "+
		"FROM %s as A "+
		"INNER JOIN %s as AA on A.id=AA.ID_Album "+
		"INNER JOIN %s as ART on ART.id=AA.ID_Artist "+
		"WHERE A.id=?", albumTable, artistAlbumTable, artistTable)
	rows, err := a.db.Query(query, idAlbum)
	if err != nil {
		return album, err
	}
	for rows.Next() {
		err = rows.Scan(&album.Id, &album.Name, &album.Cover, &date, &artistId, &artistName,
			&album.AlbumType, &album.TotalTracks, &album.SpotifyURL, &album.Popularity)
		if err != nil {
			return album, err
		}
		artist.Id = artistId
		artist.Name = artistName
		album.ReleaseDate, err = time.Parse(time.DateTime, date)
		album.Artist = append(album.Artist, artist)
		if err != nil {
			return album, err
		}
	}
	return album, err
}

func (a *AlbumMysql) SelectAllAlbumsForUser(idUser int64) ([]music.Album, error) {
	var albums []music.Album
	var album music.Album
	var date, artistID, artistName string
	var artist music.Artist
	var existInAlbums = false
	query := fmt.Sprintf("SELECT A.id,A.name,A.cover,A.release_date,ART.id,ART.name,A.album_type,A.total_tracks,A.spotify_URL,A.popularity FROM %s as A "+
		"INNER JOIN %s as UA on A.id=UA.ID_Album "+
		"INNER JOIN %s as AA on A.id=AA.ID_Album "+
		"INNER JOIN %s as ART on ART.id=AA.ID_Artist "+
		"WHERE UA.ID_User=?", albumTable, userAlbumTable, artistAlbumTable, artistTable)
	rows, err := a.db.Query(query, idUser)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&album.Id, &album.Name, &album.Cover, &date, &artistID, &artistName,
			&album.AlbumType, &album.TotalTracks, &album.SpotifyURL, &album.Popularity)
		if err != nil {
			return nil, err
		}
		artist.Id = artistID
		artist.Name = artistName
		album.ReleaseDate, err = time.Parse(time.DateTime, date)
		if err != nil {
			return nil, err
		}
		for key, al := range albums {
			if album.Id == al.Id {
				albums[key].Artist = append(albums[key].Artist, artist)
				existInAlbums = true
				break
			}
		}
		if !existInAlbums {
			album.Artist = append(album.Artist, artist)
			albums = append(albums, album)
		}
		album = music.Album{}
		existInAlbums = false
	}
	return albums, err
}

func (a *AlbumMysql) SelectAlbumsByArtist(idArtist string) ([]music.Album, error) {
	var albums []music.Album
	var album music.Album
	var date, artistID, artistName string
	var artist music.Artist
	var existInAlbums = false
	query := fmt.Sprintf("SELECT A.id,A.name,A.cover,A.release_date,ART.id,ART.name,A.album_type,A.total_tracks,A.spotify_URL,A.popularity FROM albums as A "+
		"INNER JOIN artist_album as AA on A.id=AA.ID_Album "+
		"INNER JOIN artists as ART on ART.id=AA.ID_Artist "+
		"where A.id in "+
		"(SELECT A.id FROM %s as A "+
		"INNER JOIN %s as AA on A.id=AA.ID_Album "+
		"INNER JOIN %s as ART on ART.id=AA.ID_Artist "+
		"WHERE ART.id=?)", albumTable, artistAlbumTable, artistTable)
	rows, err := a.db.Query(query, idArtist)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&album.Id, &album.Name, &album.Cover, &date, &artistID, &artistName,
			&album.AlbumType, &album.TotalTracks, &album.SpotifyURL, &album.Popularity)
		if err != nil {
			return nil, err
		}
		artist.Id = artistID
		artist.Name = artistName
		album.ReleaseDate, err = time.Parse(time.DateTime, date)
		if err != nil {
			return nil, err
		}
		for key, al := range albums {
			if album.Id == al.Id {
				albums[key].Artist = append(albums[key].Artist, artist)
				existInAlbums = true
				break
			}
		}
		if !existInAlbums {
			album.Artist = append(album.Artist, artist)
			albums = append(albums, album)
		}
		existInAlbums = false
		album = music.Album{}
	}
	return albums, err
}

func (a *AlbumMysql) DeleteAlbumFromFav(idAlbum string, idUser int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE ID_Album=? AND ID_User=?", userAlbumTable)
	_, err := a.db.Exec(query, idAlbum, idUser)
	return err
}

func (a *AlbumMysql) AddAlbumToFav(idUser int64, album music.Album) error {
	transaction, err := a.db.Begin()
	if err != nil {
		return err
	}
	err = a.InsertAlbumToDB(transaction, album)
	if err != nil {
		transaction.Rollback()
		return err
	}
	query := fmt.Sprintf("INSERT INTO %s values (?,?)", userAlbumTable)
	_, err = transaction.Exec(query, idUser, album.Id)
	if err != nil {
		transaction.Rollback()
		return err
	}
	return transaction.Commit()
}

func (a *AlbumMysql) InsertAlbumToDB(transaction *sql.Tx, album music.Album) error {
	alreadyExist, err := a.AlbumAlreadyExist(album.Id)
	if !alreadyExist {
		query := fmt.Sprintf("INSERT INTO %s values (?,?,?,?,?,?,?,?)", albumTable)
		_, err = transaction.Exec(query, album.Id, album.Name, album.Cover, album.ReleaseDate, album.AlbumType,
			album.TotalTracks, album.SpotifyURL, album.Popularity)
		if err != nil {
			transaction.Rollback()
			return err
		}
		for _, artist := range album.Artist {
			query = fmt.Sprintf("INSERT INTO %s values (?,?)", artistAlbumTable)
			_, err = transaction.Exec(query, artist.Id, album.Id)
			if err != nil {
				transaction.Rollback()
				return err
			}
		}
	}
	return err
}

func (a *AlbumMysql) StartTransaction() (*sql.Tx, error) {
	return a.db.Begin()
}

func (a *AlbumMysql) AlbumAlreadyExist(albumID string) (bool, error) {
	query := fmt.Sprintf("Select id from %s where id=?", albumTable)
	rows, err := a.db.Query(query, albumID)
	if err != nil {
		return false, err
	}
	return rows.Next(), err
}

func (a *AlbumMysql) SelectAlbumsByGenre(userId int64) ([]music.Album, error) {
	var albums []music.Album
	var album music.Album
	var date, artistID, artistName string
	var artist music.Artist
	var existInAlbums = false
	query := fmt.Sprintf("SELECT DISTINCT A.id,A.name,A.cover,A.release_date,ART.id,ART.name,A.album_type,A.total_tracks,A.spotify_URL,A.popularity FROM %s as A "+
		"INNER JOIN %s as AA on A.id=AA.ID_Album "+
		"INNER JOIN %s as ART on ART.id=AA.ID_Artist "+
		"INNER JOIN %s as AG on ART.id=AG.ID_Artist "+
		"INNER JOIN "+
		"(SELECT G.id as id FROM %s as G "+
		"INNER JOIN %s as UG ON G.id=UG.ID_Genre "+
		"WHERE UG.ID_User=?) AS GEN on GEN.id=AG.ID_Genre "+
		"WHERE A.id not in "+
		"(SELECT AL.id FROM %s AS AL "+
		"INNER JOIN %s AS UA ON UA.ID_Album=AL.id "+
		"WHERE ID_User=?) "+
		"ORDER BY A.popularity desc "+
		"LIMIT 10", albumTable, artistAlbumTable, artistTable, artistGenreTable, genresTable, userGenreTable, albumTable, userAlbumTable)
	rows, err := a.db.Query(query, userId, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&album.Id, &album.Name, &album.Cover, &date, &artistID, &artistName,
			&album.AlbumType, &album.TotalTracks, &album.SpotifyURL, &album.Popularity)
		if err != nil {
			return nil, err
		}
		artist.Id = artistID
		artist.Name = artistName
		album.ReleaseDate, err = time.Parse(time.DateTime, date)
		if err != nil {
			return nil, err
		}
		for key, al := range albums {
			if album.Id == al.Id {
				albums[key].Artist = append(albums[key].Artist, artist)
				existInAlbums = true
				break
			}
		}
		if !existInAlbums {
			album.Artist = append(album.Artist, artist)
			albums = append(albums, album)
		}
		existInAlbums = false
		album = music.Album{}
	}
	return albums, err
}

func (a *AlbumMysql) SelectAlbumDuration(albumId string) (int64, error) {
	var duration sql.NullInt64
	query := fmt.Sprintf("Select sum(T.duration) as duration from %s as T "+
		"INNER JOIN %s AS A on A.id=T.ID_Album "+
		"WHERE A.id=?", trackTable, albumTable)
	rows, err := a.db.Query(query, albumId)
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

func (a *AlbumMysql) IsAllTracksInDB(albumId string, total int64) (bool, error) {
	var count sql.NullInt64
	query := fmt.Sprintf("SELECT count(T.id) as count FROM %s AS T "+
		"INNER JOIN %s AS A on A.id=T.ID_Album "+
		"WHERE A.id=?", trackTable, albumTable)
	rows, err := a.db.Query(query, albumId)
	if err != nil {
		return false, err
	}
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return false, err
		}
	}
	if count.Valid {
		return count.Int64 == total, err
	} else {
		return false, err
	}
}
