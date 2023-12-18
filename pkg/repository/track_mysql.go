package repository

import (
	"database/sql"
	"fmt"
	music "github.com/bear1278/MusicWave"
)

type TrackMysql struct {
	db *sql.DB
}

func NewTrackMysql(db *sql.DB) *TrackMysql {
	return &TrackMysql{db: db}
}

func (t *TrackMysql) IsTrackExistInPlaylist(trackId string, playlistId int64) (bool, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE ID_Playlist=? and ID_Track=?", playlistTrackTable)
	rows, err := t.db.Query(query, playlistId, trackId)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, err
	} else {
		return false, err
	}
}

func (t *TrackMysql) AddTrack(idPlaylist int64, track music.Track) (bool, error) {
	transaction, err := t.db.Begin()
	if err != nil {
		return false, err
	}
	alreadyExist, err := t.InsertTrackToDB(transaction, track)
	if err != nil {
		transaction.Rollback()
		return false, err
	}
	query := fmt.Sprintf("Insert into %s values (?,?)", playlistTrackTable)
	_, err = transaction.Exec(query, idPlaylist, track.Id)
	if err != nil {
		transaction.Rollback()
		return false, err
	}
	return alreadyExist, transaction.Commit()
}

func (t *TrackMysql) InsertTrackToDB(transaction *sql.Tx, track music.Track) (bool, error) {
	alreadyExist, err := t.TrackAlreadyExist(track.Id)
	if err != nil {
		return false, err
	}
	if !alreadyExist {
		query := fmt.Sprintf("Insert into %s  values (?,?,?,?,?,?,?)", trackTable)
		_, err = transaction.Exec(query, track.Id, track.Name, track.Duration, track.Cover,
			track.Album.Id, track.SpotifyURL, track.Popularity)
		if err != nil {
			transaction.Rollback()
			return false, err
		}
	}
	return alreadyExist, err
}

func (t *TrackMysql) ExcludeTrack(idTrack string, idPlaylist int64) error {
	query := fmt.Sprintf("Delete from %s where ID_playlist=? and ID_Track=?", playlistTrackTable)
	_, err := t.db.Exec(query, idPlaylist, idTrack)
	return err
}

func (t *TrackMysql) SelectAllTracks(idPlaylist int64) ([]music.Track, error) {
	var tracks []music.Track
	var track music.Track
	var artist music.Artist
	var artistID, artistName string
	var existInTracks = false
	query := fmt.Sprintf("SELECT T.id,T.name,T.duration,T.cover,A.id,A.name, ART.id, ART.name,T.spotify_URL,T.popularity FROM %s as T "+
		"INNER JOIN %s as PT on T.id=PT.ID_Track "+
		"INNER JOIN %s as A ON T.ID_Album=A.id "+
		"INNER JOIN %s as AA on A.id=AA.ID_Album "+
		"INNER JOIN %s as Art on ART.id=AA.ID_Artist "+
		"WHERE PT.ID_Playlist=?", trackTable, playlistTrackTable, albumTable, artistAlbumTable, artistTable)
	rows, err := t.db.Query(query, idPlaylist)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&track.Id, &track.Name, &track.Duration, &track.Cover,
			&track.Album.Id, &track.Album.Name, &artistID, &artistName, &track.SpotifyURL, &track.Popularity)
		artist.Id = artistID
		artist.Name = artistName
		if err != nil {
			return nil, err
		}
		for key, tr := range tracks {
			if track.Id == tr.Id {
				tracks[key].Album.Artist = append(tracks[key].Album.Artist, artist)
				existInTracks = true
				break
			}
		}
		if !existInTracks {
			track.Album.Artist = append(track.Album.Artist, artist)
			tracks = append(tracks, track)
		}
		existInTracks = false
		track = music.Track{}
	}
	return tracks, err
}

func (t *TrackMysql) SelectTrackById(idTrack string) (music.Track, error) {
	var track music.Track
	var artist music.Artist
	var artistID, artistName string
	query := fmt.Sprintf("SELECT T.id,T.name,T.duration,T.cover,A.id,A.name, ART.id, ART.name,T.spotify_URL,T.popularity FROM %s as T "+
		"INNER JOIN %s as A ON T.ID_Album=A.id "+
		"INNER JOIN %s as AA on A.id=AA.ID_Album "+
		"INNER JOIN %s as Art on ART.id=AA.ID_Artist "+
		"WHERE T.id=?", trackTable, albumTable, artistAlbumTable, artistTable)
	rows, err := t.db.Query(query, idTrack)
	if err != nil {
		return track, err
	}
	for rows.Next() {
		err = rows.Scan(&track.Id, &track.Name, &track.Duration, &track.Cover,
			&track.Album.Id, &track.Album.Name, &artistID, &artistName, &track.SpotifyURL, &track.Popularity)
		if err != nil {
			return track, err
		}
		artist.Id = artistID
		artist.Name = artistName
		track.Album.Artist = append(track.Album.Artist, artist)
	}
	return track, err
}

func (t *TrackMysql) SelectTracksFromAlbum(idAlbum string) ([]music.Track, error) {
	var tracks []music.Track
	var track music.Track
	var artist music.Artist
	var artistID, artistName string
	var existInTracks = false
	query := fmt.Sprintf("SELECT T.id,T.name,T.duration,T.cover,A.id,A.name, ART.id, ART.name,T.spotify_URL,T.popularity FROM %s as T "+
		"INNER JOIN %s as A ON T.ID_Album=A.id "+
		"INNER JOIN %s as AA on A.id=AA.ID_Album "+
		"INNER JOIN %s as Art on ART.id=AA.ID_Artist "+
		"WHERE A.id=?", trackTable, albumTable, artistAlbumTable, artistTable)
	rows, err := t.db.Query(query, idAlbum)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&track.Id, &track.Name, &track.Duration, &track.Cover,
			&track.Album.Id, &track.Album.Name, &artistID, &artistName, &track.SpotifyURL, &track.Popularity)
		artist.Id = artistID
		artist.Name = artistName
		if err != nil {
			return nil, err
		}
		for key, tr := range tracks {
			if track.Id == tr.Id {
				tracks[key].Album.Artist = append(tracks[key].Album.Artist, artist)
				existInTracks = true
				break
			}
		}
		if !existInTracks {
			track.Album.Artist = append(track.Album.Artist, artist)
			tracks = append(tracks, track)
		}
		track = music.Track{}
		existInTracks = false
	}
	return tracks, err
}

func (t *TrackMysql) StartTransaction() (*sql.Tx, error) {
	return t.db.Begin()
}

func (t *TrackMysql) TrackAlreadyExist(trackID string) (bool, error) {
	query := fmt.Sprintf("Select id from %s where id=?", trackTable)
	rows, err := t.db.Query(query, trackID)
	if err != nil {
		return false, err
	}
	return rows.Next(), err
}

func (t *TrackMysql) SelectTracksByGenre(userId int64) ([]music.Track, error) {
	var track music.Track
	var tracks []music.Track
	var artist music.Artist
	var artistID, artistName string
	var existInTracks = false
	query := fmt.Sprintf("SELECT DISTINCT T.id,T.name,T.duration,T.cover,A.id,A.name, ART.id, ART.name,T.spotify_URL,T.popularity FROM %s as T "+
		"INNER JOIN %s as A on T.ID_album=A.id "+
		"INNER JOIN %s as AA on A.id=AA.ID_Album "+
		"INNER JOIN %s as ART on ART.id=AA.ID_Artist "+
		"INNER JOIN %s as AG on ART.id=AG.ID_Artist "+
		"INNER JOIN "+
		"(SELECT G.id as id FROM %s as G "+
		"INNER JOIN %s as UG ON G.id=UG.ID_Genre "+
		"WHERE UG.ID_User=?) AS GEN on GEN.id=AG.ID_Genre "+
		"WHERE T.id not in "+
		"(SELECT T.id FROM %s AS T "+
		"INNER JOIN %s as PT ON PT.ID_Track=T.id "+
		"Inner join %s as P on P.id=PT.ID_Playlist "+
		"WHERE P.author=? and P.name='favorites') "+
		"ORDER BY A.popularity desc "+
		"LIMIT 10", trackTable, albumTable, artistAlbumTable, artistTable, artistGenreTable, genresTable, userGenreTable, trackTable, playlistTrackTable, playlistTable)
	rows, err := t.db.Query(query, userId, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&track.Id, &track.Name, &track.Duration, &track.Cover,
			&track.Album.Id, &track.Album.Name, &artistID, &artistName, &track.SpotifyURL, &track.Popularity)
		artist.Id = artistID
		artist.Name = artistName
		if err != nil {
			return nil, err
		}
		for key, tr := range tracks {
			if track.Id == tr.Id {
				tracks[key].Album.Artist = append(tracks[key].Album.Artist, artist)
				existInTracks = true
				break
			}
		}
		if !existInTracks {
			track.Album.Artist = append(track.Album.Artist, artist)
			tracks = append(tracks, track)
		}
		track = music.Track{}
		existInTracks = false
	}
	return tracks, err
}

func (t *TrackMysql) SelectTopTracksOfArtist(artistId string) ([]music.Track, error) {
	var track music.Track
	var tracks []music.Track
	var artist music.Artist
	var artistID, artistName string
	var existInTracks = false
	query := fmt.Sprintf("SELECT DISTINCT T.id,T.name,T.duration,T.cover,AL.id,AL.name, ART.id, ART.name,T.spotify_URL,T.popularity FROM %s as T "+
		"INNER JOIN %s AS AL ON T.ID_Album=AL.id "+
		"INNER JOIN %s AS AA ON AL.id=AA.ID_Album "+
		"INNER JOIN %s AS ART ON ART.id=AA.ID_Artist "+
		"WHERE ART.id=? "+
		"ORDER BY T.popularity desc "+
		"LIMIT 5", trackTable, albumTable, artistAlbumTable, artistTable)
	rows, err := t.db.Query(query, artistId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&track.Id, &track.Name, &track.Duration, &track.Cover,
			&track.Album.Id, &track.Album.Name, &artistID, &artistName, &track.SpotifyURL, &track.Popularity)
		artist.Id = artistID
		artist.Name = artistName
		if err != nil {
			return nil, err
		}
		for key, tr := range tracks {
			if track.Id == tr.Id {
				tracks[key].Album.Artist = append(tracks[key].Album.Artist, artist)
				existInTracks = true
				break
			}
		}
		if !existInTracks {
			track.Album.Artist = append(track.Album.Artist, artist)
			tracks = append(tracks, track)
		}
		track = music.Track{}
		existInTracks = false
	}
	return tracks, err
}

func (t *TrackMysql) SelectMostPopularTracksFroUser(userId int64) ([]music.Track, error) {
	var track music.Track
	var tracks []music.Track
	var artist music.Artist
	var artistID, artistName string
	var existInTracks = false
	var count int64
	query := fmt.Sprintf("select id,name,duration,cover,albumid,album,artistid,artist,spotifyURL,popularity, count(playlist) as count from "+
		"((select T.id as id,T.name as name ,T.duration as duration,T.cover as cover,AL.id as albumid, "+
		"AL.name as album, ART.id as artistid, ART.name as artist, "+
		"T.spotify_URL as spotifyURL,T.popularity as popularity, PT.ID_Playlist as playlist from %s as T "+
		"INNER JOIN %s AS AL ON T.ID_Album=AL.id "+
		"INNER JOIN %s AS AA ON AL.id=AA.ID_Album "+
		"INNER JOIN %s AS ART ON ART.id=AA.ID_Artist "+
		"inner join %s as PT on T.id=PT.ID_Track "+
		"inner join %s as UP on UP.ID_Playlist=PT.ID_Playlist "+
		"where UP.ID_User=?) "+
		"union all "+
		"(select T.id as id,T.name as name ,T.duration as duration,T.cover as cover,AL.id as albumid,AL.name as album, "+
		"ART.id as artistid, ART.name as artist,T.spotify_URL as spotifyURL,T.popularity as popularity, "+
		"P.id as playlist from %s as T "+
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
	rows, err := t.db.Query(query, userId, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&track.Id, &track.Name, &track.Duration, &track.Cover,
			&track.Album.Id, &track.Album.Name, &artistID, &artistName, &track.SpotifyURL, &track.Popularity, &count)
		artist.Id = artistID
		artist.Name = artistName
		if err != nil {
			return nil, err
		}
		for key, tr := range tracks {
			if track.Id == tr.Id {
				tracks[key].Album.Artist = append(tracks[key].Album.Artist, artist)
				existInTracks = true
				break
			}
		}
		if !existInTracks {
			track.Album.Artist = append(track.Album.Artist, artist)
			tracks = append(tracks, track)
		}
		track = music.Track{}
		existInTracks = false
	}
	return tracks, err
}
