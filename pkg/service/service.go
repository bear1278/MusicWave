package service

import (
	music "github.com/bear1278/MusicWave"
	"github.com/bear1278/MusicWave/pkg/repository"
	"mime/multipart"
)

type Authorization interface {
	CreateUser(user music.User) (int64, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int64, error)
	FillHtml() ([]music.Genre, error)
	InsertRecommendation(genres []music.Genre, userId int64) error
	GeneratePasswordHash(password string) string
	GenerateTokenForReset(username, email string) (music.User, string, error)
	SetNewPassword(id int64, password string) error
}

type PlaylistService interface {
	NewPlaylist(playlist music.Playlist, userId int64) (int64, error)
	DeletePlaylist(playlistID, userID int64) error
	GetAllPlaylist(idUser int64) ([]music.Playlist, []music.Playlist, error)
	GetById(playlist int64) (music.Playlist, error)
	UpdatePlaylist(playlist music.Playlist) error
	ExcludePlaylist(idPlaylist, idUser int64) error
	AddPlaylist(idPlaylist, idUser int64) error
	SearchPlaylist(search string) ([]music.Playlist, error)
	GetAllPlaylistForAdmin() ([]music.Playlist, error)
	DeletePlaylistByAdmin(playlistID int64, reason string) error
	GetUserFavorites(userID int64) (music.Playlist, error)
	GerPlaylistsRec(userId int64) ([]music.Playlist, error)
	IsAddedToSpotify(playlistId int64) (bool, error)
	UpdateIsAdded(playlistId int64) error
	GetDurationOfAllPlaylists(userId int64) (int64, error)
}

type TrackService interface {
	AddTrack(idPlaylist int64, track music.Track) (bool, error)
	ExcludeTrack(idPlaylist int64, idTrack string) error
	GetAllTracks(idPlaylist int64) ([]music.Track, error)
	GetTrackById(idTrack string) (music.Track, error)
	GetTrackFromAlbum(idAlbum string) ([]music.Track, error)
	AddTrackToDB(track music.Track) (bool, error)
	GetTracksForRec(userId int64) ([]music.Track, error)
	GetTopTracksOfArtist(artistId string) ([]music.Track, error)
	GetMostPopularTracksForUser(userid int64) ([]music.Track, error)
}

type AlbumService interface {
	GetAlbumById(idAlbum string) (music.Album, error)
	GetAllAlbumsForUser(idUser int64) ([]music.Album, error)
	GetAlbumsByArtist(idArtist string) ([]music.Album, error)
	ExcludeAlbum(idUser int64, idAlbum string) error
	AddAlbumToFav(idUser int64, album music.Album) error
	AddAlbumToDB(album music.Album) error
	GetAlbumsRec(userId int64) ([]music.Album, error)
}

type ArtistService interface {
	GetAllAlbumsForUser(idUser int64) ([]music.Artist, error)
	GetArtistById(idArtist string) (music.Artist, error)
	ExcludeArtist(idUser int64, idArtist string) error
	AddArtistToFav(idUser int64, artist music.Artist) error
	AddArtistToDB(artist music.Artist) error
	GetAllArtists() ([]music.Artist, error)
	DeleteArtist(artistID, reason string) error
	GetArtistsRec(userId int64) ([]music.Artist, error)
	GetMostPopularTrackForUser(userId int64) ([]music.Artist, error)
}

type GenreService interface {
	AddGenreToArtist(artistID string, genres []music.Genre) error
	GerUserGenres(userId int64) ([]music.Genre, error)
}

type UserService interface {
	ChangeUsername(userID int64, username string) error
	ChangePassword(userId int64, newPassword, oldPassword string) error
	ChangePicture(userID int64, picture string) error
	ChangeEmail(userID int64, email string) error
	GetAllUsers() ([]music.User, error)
	DeleteUser(idUser int64, reason string) error
	GetUserByID(idUser int64) (music.User, error)
}

type AdminService interface {
	CheckAdmin(userID int64) (bool, error)
	GetHistory() ([]music.History, error)
	GetReportPDF() (string, error)
	GetReportExcel() (string, error)
	GetDBInJSON() (string, error)
	ImportDBInJSON(file *multipart.FileHeader) error
	GetGenrePopularity() ([]music.Genre, error)
	GetArtistPopularity() ([]music.Artist, error)
	GetGenreDiversity() ([]music.Genre, error)
}

type Service struct {
	Authorization
	PlaylistService
	TrackService
	AlbumService
	ArtistService
	GenreService
	UserService
	AdminService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization:   NewAuthService(repo.Authorization),
		PlaylistService: NewPlaylistServiceImpl(repo.PlaylistRepo),
		TrackService:    NewTrackServiceImpl(repo.TrackRepo),
		AlbumService:    NewAlbumServiceImpl(repo.AlbumRepo),
		ArtistService:   NewArtistServiceImpl(repo.ArtistRepo),
		GenreService:    NewGenreServiceImpl(repo.GenreRepo),
		UserService:     NewUserServiceImpl(repo.UserRepo),
		AdminService:    NewAdminServiceImpl(repo.AdminRepo),
	}
}
