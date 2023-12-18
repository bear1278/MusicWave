package repository

import (
	"database/sql"
	"fmt"
	music "github.com/bear1278/MusicWave"
	"time"
)

type ArtistMysql struct {
	db *sql.DB
}

func NewArtistMysql(db *sql.DB) *ArtistMysql {
	return &ArtistMysql{db: db}
}

func (a *ArtistMysql) IsArtistExistInUserLibrary(artistId string, userId int64) (bool, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE ID_Artist=? and ID_User=?", userArtistTable)
	rows, err := a.db.Query(query, artistId, userId)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, err
	} else {
		return false, err
	}
}

func (a *ArtistMysql) SelectAllArtistForUser(idUser int64) ([]music.Artist, error) {
	var artists []music.Artist
	var artist music.Artist
	query := fmt.Sprintf("SELECT A.id,A.name,A.picture,A.spotify_URL,A.popularity FROM %s as A "+
		"INNER JOIN %s as UA on A.id=UA.ID_Artist "+
		"WHERE UA.ID_User=?", artistTable, userArtistTable)
	rows, err := a.db.Query(query, idUser)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&artist.Id, &artist.Name, &artist.Picture, &artist.SpotifyURL, &artist.Popularity)
		if err != nil {
			return nil, err
		}
		artists = append(artists, artist)
	}
	return artists, err
}

func (a *ArtistMysql) SelectArtistById(idArtist string) (music.Artist, error) {
	var artist music.Artist
	query := fmt.Sprintf("Select id,name,picture,spotify_URL,popularity FROM %s WHERE id=?", artistTable)
	rows, err := a.db.Query(query, idArtist)
	if err != nil {
		return music.Artist{}, err
	}
	for rows.Next() {
		err = rows.Scan(&artist.Id, &artist.Name, &artist.Picture, &artist.SpotifyURL, &artist.Popularity)
		if err != nil {
			return music.Artist{}, err
		}
	}
	return artist, err
}

func (a *ArtistMysql) DeleteArtistFromFav(idArtist string, idUser int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE ID_User=? AND ID_Artist=?", userArtistTable)
	_, err := a.db.Exec(query, idUser, idArtist)
	return err
}

func (a *ArtistMysql) AddArtistToFav(idUser int64, artist music.Artist) error {
	transaction, err := a.db.Begin()
	if err != nil {
		return err
	}
	err = a.InsertArtistToDB(transaction, artist)
	if err != nil {
		transaction.Rollback()
		return err
	}
	query := fmt.Sprintf("INSERT INTO %s values (?,?)", userArtistTable)
	_, err = transaction.Exec(query, idUser, artist.Id)
	if err != nil {
		transaction.Rollback()
		return err
	}
	return transaction.Commit()
}

func (a *ArtistMysql) InsertArtistToDB(transaction *sql.Tx, artist music.Artist) error {
	alreadyExist, err := a.ArtistAlreadyExist(artist.Id)
	if !alreadyExist {
		query := fmt.Sprintf("INSERT INTO %s values (?,?,?,?,?,?)", artistTable)
		_, err = transaction.Exec(query, artist.Id, artist.Name, artist.Picture, artist.SpotifyURL, artist.Popularity, time.Now())
		if err != nil {
			transaction.Rollback()
			return err
		}
	}
	return err
}

func (a *ArtistMysql) StartTransaction() (*sql.Tx, error) {
	return a.db.Begin()
}

func (a *ArtistMysql) SelectAllArtists() ([]music.Artist, error) {
	var artists []music.Artist
	var artist music.Artist
	var date string
	query := fmt.Sprintf("SELECT * FROM %s", artistTable)
	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&artist.Id, &artist.Name, &artist.Picture, &artist.SpotifyURL, &artist.Popularity, &date)
		if err != nil {
			return nil, err
		}
		artist.AddedDate, err = time.Parse(time.DateTime, date)
		if err != nil {
			return nil, err
		}
		artists = append(artists, artist)
		artist = music.Artist{}
	}
	return artists, err
}

func (a *ArtistMysql) DeleteArtist(artistID string, reason string) error {
	var name string
	transaction, err := a.db.Begin()
	if err != nil {
		return err
	}
	query := fmt.Sprintf("SELECT name FROM %s WHERE id=?", artistTable)
	rows, err := transaction.Query(query, artistID)
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
	_, err = transaction.Exec(query, name, "artist", reason, time.Now())
	if err != nil {
		transaction.Rollback()
		return err
	}
	query = fmt.Sprintf("DELETE FROM %s "+
		"WHERE id in "+
		"(SELECT A.id from (SELECT * FROM %s) as A "+
		"INNER JOIN %s as AA on A.id=AA.ID_Album "+
		"WHERE AA.ID_Artist=?)", albumTable, albumTable, artistAlbumTable)
	_, err = transaction.Exec(query, artistID)
	if err != nil {
		transaction.Rollback()
		return err
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE id=?", artistTable)
	_, err = transaction.Exec(query, artistID)
	if err != nil {
		transaction.Rollback()
		return err
	}
	return transaction.Commit()
}

func (a *ArtistMysql) ArtistAlreadyExist(artistID string) (bool, error) {
	query := fmt.Sprintf("Select id from %s where id=?", artistTable)
	rows, err := a.db.Query(query, artistID)
	if err != nil {
		return false, err
	}
	return rows.Next(), err
}

func (a *ArtistMysql) SelectArtistsByGenre(userId int64) ([]music.Artist, error) {
	var artists []music.Artist
	var artist music.Artist
	query := fmt.Sprintf("SELECT DISTINCT A.id,A.name,A.picture,A.spotify_URL,A.popularity FROM %s as A "+
		"INNER JOIN %s as AG on A.id=AG.ID_Artist "+
		"INNER JOIN "+
		"(SELECT G.id as id FROM %s as G "+
		"INNER JOIN %s as UG ON G.id=UG.ID_Genre "+
		"WHERE UG.ID_User=?) AS GEN on GEN.id=AG.ID_Genre "+
		"WHERE A.id not in "+
		"(SELECT AR.id FROM %s AS AR "+
		"INNER JOIN %s AS UA ON UA.ID_Artist=AR.id "+
		"WHERE ID_User=?) "+
		"ORDER BY A.popularity desc "+
		"LIMIT 10", artistTable, artistGenreTable, genresTable, userGenreTable, artistTable, userArtistTable)
	rows, err := a.db.Query(query, userId, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&artist.Id, &artist.Name, &artist.Picture, &artist.SpotifyURL, &artist.Popularity)
		if err != nil {
			return nil, err
		}
		artists = append(artists, artist)
	}
	return artists, err
}

func (a *ArtistMysql) SelectMostPopularArtistForUser(userID int64) ([]music.Artist, error) {
	var artists []music.Artist
	var artist music.Artist
	var count int64
	query := fmt.Sprintf("select id,name,picture,spotifyURL,popularity, count(playlist) as count from "+
		"((select ART.id as id,ART.name as name ,ART.picture as picture, "+
		"ART.spotify_URL as spotifyURL,ART.popularity as popularity, PT.ID_Playlist as playlist from %s as T "+
		"INNER JOIN %s AS AL ON T.ID_Album=AL.id "+
		"INNER JOIN %s AS AA ON AL.id=AA.ID_Album "+
		"INNER JOIN %s AS ART ON ART.id=AA.ID_Artist "+
		"inner join %s as PT on T.id=PT.ID_Track "+
		"inner join %s as UP on UP.ID_Playlist=PT.ID_Playlist "+
		"where UP.ID_User=?) "+
		"union all "+
		"(select ART.id as id,ART.name as name ,ART.picture as picture, "+
		"ART.spotify_URL as spotifyURL,ART.popularity as popularity, PT.ID_Playlist as playlist from %s as T "+
		"INNER JOIN %s AS AL ON T.ID_Album=AL.id "+
		"INNER JOIN %s AS AA ON AL.id=AA.ID_Album "+
		"INNER JOIN %s AS ART ON ART.id=AA.ID_Artist "+
		"inner join %s as PT on T.id=PT.ID_Track "+
		"inner join %s as P on P.id=PT.ID_Playlist "+
		"where P.author=? and P.name='favorites')) as info "+
		"group by id "+
		"order by count desc "+
		"limit 5", trackTable, albumTable, artistAlbumTable, artistTable, playlistTrackTable,
		userPlaylistTable, trackTable, albumTable, artistAlbumTable, artistTable, playlistTrackTable, playlistTable)
	rows, err := a.db.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&artist.Id, &artist.Name, &artist.Picture, &artist.SpotifyURL, &artist.Popularity, &count)
		if err != nil {
			return nil, err
		}
		artists = append(artists, artist)
	}
	return artists, err
}
