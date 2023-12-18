package repository

import (
	"database/sql"
	music "github.com/bear1278/MusicWave"
)

type Authorization interface {
	CreateUser(user music.User) (int64, error)
	GetUser(username, password string) (music.User, error)
	GetAllGenres() ([]music.Genre, error)
	InsertUserGenre(genres []music.Genre, userId int64) error
	SelectUsersEmail(username, email string) (music.User, error)
	SetNewPassword(id int64, password string) error
}

type PlaylistRepo interface {
	CreatePlaylist(playlist music.Playlist) (int64, error)
	DeletePlaylist(playlistID, userID int64) error
	SelectAllPlaylists(idUser int64) ([]music.Playlist, error)
	SelectById(playlist int64) (music.Playlist, error)
	UpdatePlaylist(playlist music.Playlist) error
	ExcludePlaylist(idPlaylist, idUser int64) error
	AddPlaylist(idPlaylist, idUser int64) error
	SelectPlaylistForSearch(search string) ([]music.Playlist, error)
	SelectAllPlaylistsForAdmin() ([]music.Playlist, error)
	DeletePlaylistByAdmin(playlistId int64, reason string) error
	GetUserFavourites(userID int64) (music.Playlist, error)
	SelectPlaylistsRec(userId int64) ([]music.Playlist, error)
	SelectPlaylistDuration(playlistId int64) (int64, error)
	SelectIsAddedToSpotify(playlistId int64) (bool, error)
	UpdateIsAdded(playlistId int64) error
	IsExistInUserLibrary(playlistId, userId int64) (bool, error)
}

type TrackRepo interface {
	AddTrack(idPlaylist int64, track music.Track) (bool, error)
	ExcludeTrack(idTrack string, idPlaylist int64) error
	SelectAllTracks(idPlaylist int64) ([]music.Track, error)
	SelectTrackById(idTrack string) (music.Track, error)
	SelectTracksFromAlbum(idAlbum string) ([]music.Track, error)
	InsertTrackToDB(transaction *sql.Tx, track music.Track) (bool, error)
	StartTransaction() (*sql.Tx, error)
	SelectTracksByGenre(userId int64) ([]music.Track, error)
	SelectTopTracksOfArtist(artistId string) ([]music.Track, error)
	SelectMostPopularTracksFroUser(userId int64) ([]music.Track, error)
	IsTrackExistInPlaylist(trackId string, playlistId int64) (bool, error)
}

type AlbumRepo interface {
	GetAlbumById(idAlbum string) (music.Album, error)
	SelectAllAlbumsForUser(idUser int64) ([]music.Album, error)
	SelectAlbumsByArtist(idArtist string) ([]music.Album, error)
	DeleteAlbumFromFav(idAlbum string, idUser int64) error
	AddAlbumToFav(idUser int64, album music.Album) error
	InsertAlbumToDB(transaction *sql.Tx, album music.Album) error
	StartTransaction() (*sql.Tx, error)
	SelectAlbumsByGenre(userId int64) ([]music.Album, error)
	SelectAlbumDuration(albumId string) (int64, error)
	IsAllTracksInDB(albumId string, total int64) (bool, error)
	IsAlbumExistInUserLibrary(albumId string, userId int64) (bool, error)
}

type ArtistRepo interface {
	SelectAllArtistForUser(idUser int64) ([]music.Artist, error)
	SelectArtistById(idArtist string) (music.Artist, error)
	DeleteArtistFromFav(idArtist string, idUser int64) error
	AddArtistToFav(idUser int64, artist music.Artist) error
	InsertArtistToDB(transaction *sql.Tx, artist music.Artist) error
	StartTransaction() (*sql.Tx, error)
	SelectAllArtists() ([]music.Artist, error)
	DeleteArtist(artistID string, reason string) error
	SelectArtistsByGenre(userId int64) ([]music.Artist, error)
	SelectMostPopularArtistForUser(userID int64) ([]music.Artist, error)
	IsArtistExistInUserLibrary(artistId string, userId int64) (bool, error)
}

type GenreRepo interface {
	InsertGenre(artistID string, genres []music.Genre) error
	SelectUserGenre(userID int64) ([]music.Genre, error)
}

type UserRepo interface {
	UpdateUsername(userID int64, username string) error
	UpdatePassword(userId int64, password string) error
	SelectPassword(userId int64) (string, error)
	UpdatePicture(userId int64, picture string) error
	UpdateEmail(userId int64, email string) error
	SelectAllUsers() ([]music.User, error)
	DeleteUser(idUser int64, reason string) error
	SelectUserById(idUser int64) (music.User, error)
}

type AdminRepo interface {
	SelectArtistPopularity() ([]music.Artist, error)
	SelectGenrePopularity() ([]music.Genre, error)
	SelectGenreDiversity() ([]music.Genre, error)
	SelectAdminId() (int64, error)
	SelectHistory() ([]music.History, error)
	SelectTable(tableName string) ([]map[string]interface{}, error)
	InsertArtist(object music.Artist) error
	InsertUser(user music.User) error
	InsertPlaylist(playlist music.Playlist) error
	InsertAlbum(album music.Album) error
	InsertTrack(object music.Track) error
	InsertGenre(object music.Genre) error
	InsertHistory(object music.History) error
	InsertArtistAlbum(object music.ArtistAlbum) error
	InsertUserAlbum(object music.UserAlbum) error
	InsertUserArtist(object music.UserArtist) error
	InsertUserGenre(object music.UserGenre) error
	InsertUserPlaylist(object music.UserPlaylist) error
	InsertPlaylistTrack(object music.PlaylistTrack) error
	InsertArtistGenre(object music.ArtistGenre) error
}

type Repository struct {
	Authorization
	PlaylistRepo
	TrackRepo
	AlbumRepo
	ArtistRepo
	GenreRepo
	UserRepo
	AdminRepo
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthMysql(db),
		PlaylistRepo:  NewPlaylistMysql(db),
		TrackRepo:     NewTrackMysql(db),
		AlbumRepo:     NewAlbumMysql(db),
		ArtistRepo:    NewArtistMysql(db),
		GenreRepo:     NewGenreMysql(db),
		UserRepo:      NewUserMysql(db),
		AdminRepo:     NewAdminMysql(db),
	}
}
