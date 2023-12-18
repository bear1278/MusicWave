package repository

import (
	"database/sql"
	"fmt"
	music "github.com/bear1278/MusicWave"
	"time"
)

type AdminMysql struct {
	db *sql.DB
}

func (a *AdminMysql) SelectAdminId() (int64, error) {
	var id int64
	query := fmt.Sprintf("SELECT id FROM %s WHERE username='admin'", usersTable)
	rows, err := a.db.Query(query)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return 0, err
		}
	}
	return id, err
}

func NewAdminMysql(db *sql.DB) *AdminMysql {
	return &AdminMysql{db: db}
}

func (a *AdminMysql) SelectHistory() ([]music.History, error) {
	var history music.History
	var histories []music.History
	var date string
	query := fmt.Sprintf("SELECT * FROM %s", adminHistory)
	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&history.ID, &history.Name, &history.Type, &history.Reason, &date)
		if err != nil {
			return nil, err
		}
		history.DeletedDate, err = time.Parse(time.DateTime, date)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
		history = music.History{}
	}
	return histories, err
}

func (a *AdminMysql) SelectTable(tableName string) ([]map[string]interface{}, error) {
	rows, err := a.db.Query("SELECT * FROM " + tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Преобразование данных в формат JSON
	var result []map[string]interface{}
	columns, _ := rows.Columns()
	for rows.Next() {
		row := make(map[string]interface{})
		values := make([]interface{}, len(columns))
		for i := range columns {
			values[i] = new(interface{})
		}
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}
		for i, column := range columns {
			row[column] = *(values[i].(*interface{}))
		}
		result = append(result, row)
	}

	return result, err
}

func (a *AdminMysql) InsertUser(user music.User) error {
	query := fmt.Sprintf("INSERT INTO %s "+
		"(name,username,email,password,picture,create_date,modified_date) "+
		"values (?,?,?,?,?,?,?)", usersTable)
	_, err := a.db.Exec(query, user.Name, user.UserName, user.Email, user.Password, user.Picture, user.CreateDate, user.ModifiedDate)
	return err
}

func (a *AdminMysql) InsertPlaylist(playlist music.Playlist) error {
	query := fmt.Sprintf("INSERT INTO %s "+
		"(name,duration,cover,author,release_date,type,modified_date) "+
		"values (?,?,?,?,?,?,?)", playlistTable)
	_, err := a.db.Exec(query, playlist.Name, playlist.Duration, playlist.Cover,
		playlist.Author.Id, playlist.ReleaseDate, playlist.Type, playlist.ModifiedDate)
	return err
}

func (a *AdminMysql) InsertAlbum(album music.Album) error {
	query := fmt.Sprintf("INSERT INTO %s values (?,?,?,?,?,?,?,?)", albumTable)
	_, err := a.db.Exec(query, album.Id, album.Name, album.Cover, album.ReleaseDate,
		album.AlbumType, album.TotalTracks, album.SpotifyURL, album.Popularity)
	return err
}

func (a *AdminMysql) InsertArtist(object music.Artist) error {
	query := fmt.Sprintf("INSERT INTO %s values (?,?,?,?,?,?)", artistTable)
	_, err := a.db.Exec(query, object.Id, object.Name, object.Picture,
		object.SpotifyURL, object.Popularity, object.AddedDate)
	return err
}

func (a *AdminMysql) InsertTrack(object music.Track) error {
	query := fmt.Sprintf("INSERT INTO %s values (?,?,?,?,?,?,?)", artistTable)
	_, err := a.db.Exec(query, object.Id, object.Name, object.Duration, object.Cover, object.Album.Id,
		object.SpotifyURL, object.Popularity)
	return err
}

func (a *AdminMysql) InsertGenre(object music.Genre) error {
	query := fmt.Sprintf("INSERT INTO %s (name) values (?)", genresTable)
	_, err := a.db.Exec(query, object.Name)
	return err
}

func (a *AdminMysql) InsertHistory(object music.History) error {
	query := fmt.Sprintf("INSERT INTO %s (name,type,reason,deleted_date) values (?,?,?,?)", adminHistory)
	_, err := a.db.Exec(query, object.Name, object.Type, object.Reason, object.DeletedDate)
	return err
}

func (a *AdminMysql) InsertArtistAlbum(object music.ArtistAlbum) error {
	query := fmt.Sprintf("INSERT INTO %s values (?,?)", artistAlbumTable)
	_, err := a.db.Exec(query, object.IdArtist, object.IdAlbum)
	return err
}

func (a *AdminMysql) InsertUserAlbum(object music.UserAlbum) error {
	query := fmt.Sprintf("INSERT INTO %s values (?,?)", userAlbumTable)
	_, err := a.db.Exec(query, object.IdUser, object.IdAlbum)
	return err
}

func (a *AdminMysql) InsertUserArtist(object music.UserArtist) error {
	query := fmt.Sprintf("INSERT INTO %s values (?,?)", userAlbumTable)
	_, err := a.db.Exec(query, object.IdUser, object.IdArtist)
	return err
}

func (a *AdminMysql) InsertUserGenre(object music.UserGenre) error {
	query := fmt.Sprintf("INSERT INTO %s values (?,?)", userGenreTable)
	_, err := a.db.Exec(query, object.IdUser, object.IdGenre)
	return err
}

func (a *AdminMysql) InsertUserPlaylist(object music.UserPlaylist) error {
	query := fmt.Sprintf("INSERT INTO %s values (?,?)", userPlaylistTable)
	_, err := a.db.Exec(query, object.IdUser, object.IdPlaylist)
	return err
}

func (a *AdminMysql) InsertPlaylistTrack(object music.PlaylistTrack) error {
	query := fmt.Sprintf("INSERT INTO %s values (?,?)", playlistTrackTable)
	_, err := a.db.Exec(query, object.IdPlaylist, object.IdTrack)
	return err
}

func (a *AdminMysql) InsertArtistGenre(object music.ArtistGenre) error {
	query := fmt.Sprintf("INSERT INTO %s values (?,?)", artistGenreTable)
	_, err := a.db.Exec(query, object.IdArtist, object.IdGenre)
	return err
}

func (a *AdminMysql) SelectGenrePopularity() ([]music.Genre, error) {
	var genre music.Genre
	var genres []music.Genre
	query := fmt.Sprintf("SELECT G.id,G.name, count(UG.ID_user) as popularity from %s as G "+
		"INNER JOIN %s as UG on UG.ID_GENRE=G.id "+
		"GROUP BY G.id "+
		"ORDER BY popularity desc "+
		"LIMIT 10", genresTable, userGenreTable)
	rows, err := a.db.Query(query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&genre.Id, &genre.Name, &genre.Popularity)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
		genre = music.Genre{}
	}
	return genres, err
}

func (a *AdminMysql) SelectArtistPopularity() ([]music.Artist, error) {
	var artist music.Artist
	var artists []music.Artist
	query := fmt.Sprintf("SELECT A.name,popularity from %s as A "+
		"ORDER BY popularity desc "+
		"LIMIT 10", artistTable)
	rows, err := a.db.Query(query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&artist.Name, &artist.Popularity)
		if err != nil {
			return nil, err
		}
		artists = append(artists, artist)
		artist = music.Artist{}
	}
	return artists, err
}

func (a *AdminMysql) SelectGenreDiversity() ([]music.Genre, error) {
	var genre music.Genre
	var genres []music.Genre
	query := fmt.Sprintf("SELECT G.id,G.name, count(T.id) as popularity from %s as G "+
		"INNER JOIN %s as AG on AG.ID_GENRE=G.id "+
		"INNER JOIN %s as A on A.id=AG.ID_Artist "+
		"INNER JOIN %s as AA on AA.ID_Artist=A.id "+
		"INNER JOIN %s as AL on AL.id=AA.ID_Album "+
		"INNER JOIN %s as T on AL.id=T.ID_Album "+
		"GROUP BY G.id "+
		"ORDER BY popularity desc ", genresTable, artistGenreTable, artistTable, artistAlbumTable, albumTable, trackTable)
	rows, err := a.db.Query(query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&genre.Id, &genre.Name, &genre.Diversity)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
		genre = music.Genre{}
	}
	return genres, err
}
