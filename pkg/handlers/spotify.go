package handlers

import (
	"errors"
	music "github.com/bear1278/MusicWave"
	"github.com/bear1278/MusicWave/configs"
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify"
	"net/http"
	"strconv"
)

var (
	redirectURI = "http://localhost:8000/relocate" // Set this to the redirect URI you specified in the Spotify Developer Dashboard
	scopes      = []string{spotify.ScopeUserReadPrivate, spotify.ScopeUserReadEmail, spotify.ScopePlaylistModifyPublic, spotify.ScopePlaylistModifyPrivate, spotify.ScopeUserLibraryModify, spotify.ScopeUserTopRead, spotify.ScopeUserLibraryRead}
)

var (
	authAPI = spotify.NewAuthenticator(redirectURI, scopes...)
)

func (h *Handler) LoginHandler(ctx *gin.Context) {
	cfg, err := configs.Init()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	authAPI.SetAuthInfo(cfg.Spotify.ClientID, cfg.Spotify.ClientSecret)

	url := authAPI.AuthURL(cfg.Spotify.State)
	ctx.Redirect(http.StatusFound, url)
}

func (h *Handler) GetClient(ctx *gin.Context) {
	cfg, err := configs.Init()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := authAPI.TokenWithOpts(cfg.Spotify.State, ctx.Request)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, token)
}

func (h *Handler) SearchTrack(ctx *gin.Context) {
	var tracks []music.Track
	var track music.Track
	searchString := ctx.Param("string")
	var next = "/search/track/" + searchString + "/"
	var previous = "/search/track/" + searchString + "/"
	numberOfPage, err := strconv.Atoi(ctx.Param("page"))
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := GetSpotifyToken(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	client := authAPI.NewClient(token)
	items, err := client.Search(searchString, spotify.SearchTypeTrack)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	for i := 0; i < numberOfPage; i++ {
		err = client.NextTrackResults(items)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}
	for _, trackAPI := range items.Tracks.Tracks {
		track.NewTrack(trackAPI)
		tracks = append(tracks, track)
		track = music.Track{}
	}
	if items.Tracks.Next != "" {
		next += strconv.Itoa(numberOfPage + 1)
	} else {
		next = ""
	}
	if items.Tracks.Previous != "" {
		previous += strconv.Itoa(numberOfPage - 1)
	} else {
		previous = ""
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data":     tracks,
		"previous": previous,
		"nextPage": next})
}

func (h *Handler) SearchAlbum(ctx *gin.Context) {
	var albums []music.Album
	var album music.Album
	searchString := ctx.Param("string")
	var next = "/search/album/" + searchString + "/"
	var previous = "/search/album/" + searchString + "/"
	numberOfPage, err := strconv.Atoi(ctx.Param("page"))
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	token, err := GetSpotifyToken(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	client := authAPI.NewClient(token)
	items, err := client.Search(searchString, spotify.SearchTypeAlbum)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	for i := 0; i < numberOfPage; i++ {
		err = client.NextAlbumResults(items)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}
	for _, albumAPI := range items.Albums.Albums {
		album.NewAlbum(albumAPI)
		albums = append(albums, album)
		album = music.Album{}
	}
	if items.Albums.Next != "" {
		next += strconv.Itoa(numberOfPage + 1)
	} else {
		next = ""
	}
	if items.Albums.Previous != "" {
		previous += strconv.Itoa(numberOfPage - 1)
	} else {
		previous = ""
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data":     albums,
		"previous": previous,
		"nextPage": next})
}

func (h *Handler) SearchArtist(ctx *gin.Context) {
	var artists []music.Artist
	var artist music.Artist
	searchString := ctx.Param("string")
	var next = "/search/artist/" + searchString + "/"
	var previous = "/search/artist/" + searchString + "/"
	numberOfPage, err := strconv.Atoi(ctx.Param("page"))
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	token, err := GetSpotifyToken(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	client := authAPI.NewClient(token)
	items, err := client.Search(searchString, spotify.SearchTypeArtist)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	for i := 0; i < numberOfPage; i++ {
		err = client.NextArtistResults(items)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}
	for _, artistAPI := range items.Artists.Artists {
		artist.NewArtist(artistAPI)
		artists = append(artists, artist)
		artist = music.Artist{}
	}
	if items.Artists.Next != "" {
		next += strconv.Itoa(numberOfPage + 1)
	} else {
		next = ""
	}
	if items.Artists.Previous != "" {
		previous += strconv.Itoa(numberOfPage - 1)
	} else {
		previous = ""
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data":     artists,
		"previous": previous,
		"nextPage": next})
}

func (h *Handler) GetArtistByIDFromApi(artistID string, ctx *gin.Context) (music.Artist, error) {
	var artist music.Artist
	token, err := GetSpotifyToken(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return music.Artist{}, err
	}
	client := authAPI.NewClient(token)
	artistAPI, err := client.GetArtist(spotify.ID(artistID))
	if err != nil {
		return music.Artist{}, err
	}
	artist.NewArtist(*artistAPI)
	return artist, err
}

func (h *Handler) GetAlbumByIDFromApi(albumID string, ctx *gin.Context) (music.Album, error) {
	var album music.Album
	token, err := GetSpotifyToken(ctx)
	if err != nil {
		return music.Album{}, err
	}
	client := authAPI.NewClient(token)
	albumAPI, err := client.GetAlbum(spotify.ID(albumID))
	if err != nil {
		return music.Album{}, err
	}
	album.NewAlbum(albumAPI.SimpleAlbum)
	album.Popularity = int64(albumAPI.Popularity)
	album.TotalTracks = int64(albumAPI.Tracks.Total)
	return album, err
}

func (h *Handler) GetTrackByIDFromAPI(trackID string, ctx *gin.Context) (music.Track, error) {
	var track music.Track
	token, err := GetSpotifyToken(ctx)
	if err != nil {
		return music.Track{}, err
	}
	client := authAPI.NewClient(token)
	trackAPI, err := client.GetTrack(spotify.ID(trackID))
	if err != nil {
		return music.Track{}, err
	}
	track.NewTrack(*trackAPI)
	return track, err
}

func (h *Handler) GetTracksFromAlbumFromApi(albumID string, ctx *gin.Context) ([]music.Track, error) {
	var track music.Track
	var tracks []music.Track
	token, err := GetSpotifyToken(ctx)
	if err != nil {
		return nil, err
	}
	client := authAPI.NewClient(token)
	tracksAPI, err := client.GetAlbumTracks(spotify.ID(albumID))
	if err != nil {
		return nil, err
	}
	for _, tr := range tracksAPI.Tracks {
		track, err = h.GetTrackByIDFromAPI(tr.ID.String(), ctx)
		tracks = append(tracks, track)
		track = music.Track{}
	}
	return tracks, err
}

func (h *Handler) GetTop5TracksOfArtist(artistID string, ctx *gin.Context) ([]music.Track, error) {
	var track music.Track
	var tracks []music.Track
	token, err := GetSpotifyToken(ctx)
	if err != nil {
		return nil, err
	}
	client := authAPI.NewClient(token)
	tracksAPI, err := client.GetArtistsTopTracks(spotify.ID(artistID), "BY")
	if err != nil {
		return nil, err
	}
	for _, trackAPI := range tracksAPI {
		track.NewTrack(trackAPI)
		tracks = append(tracks, track)
		track = music.Track{}
	}
	return tracks, err
}

func (h *Handler) GetAlbumsOfArtistFromApi(artistID string, ctx *gin.Context) ([]music.Album, error) {
	var album music.Album
	var albums []music.Album
	token, err := GetSpotifyToken(ctx)
	if err != nil {
		return nil, err
	}
	client := authAPI.NewClient(token)
	albumsAPI, err := client.GetArtistAlbums(spotify.ID(artistID))
	if err != nil {
		return nil, err
	}
	for _, albumAPI := range albumsAPI.Albums {
		album.NewAlbum(albumAPI)
		albums = append(albums, album)
		album = music.Album{}
	}
	return albums, err
}

func (h *Handler) GetAlbumDurationFromApi(albumID string, ctx *gin.Context) (int64, error) {
	var duration int64 = 0
	token, err := GetSpotifyToken(ctx)
	if err != nil {
		return 0, err
	}
	client := authAPI.NewClient(token)
	tracksAPI, err := client.GetAlbumTracks(spotify.ID(albumID))
	if err != nil {
		return 0, err
	}
	for _, tr := range tracksAPI.Tracks {
		duration += int64(tr.Duration)
	}
	return duration, err
}

func (h *Handler) GetSearchPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "search.html", nil)
}

func (h *Handler) CreatePlaylistInSpotify(ctx *gin.Context) {
	var public bool
	var trackIds []spotify.ID
	playlistID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	userID, err := GetUserId(ctx)
	if err != nil {
		return
	}

	playlist, err := h.service.PlaylistService.GetById(playlistID)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if playlist.Author.Id != userID {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("you are not owner of this playlist").Error())
		return
	}

	if playlist.Name == "favorites" {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("you can't add favorites").Error())
		return
	}

	token, err := GetSpotifyToken(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	client := authAPI.NewClient(token)
	user, err := client.CurrentUser()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if playlist.Type == "public" {
		public = true
	} else {
		public = false
	}
	playlistAPI, err := client.CreatePlaylistForUser(user.ID, playlist.Name, "", public)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	tracks, err := h.service.TrackService.GetAllTracks(playlist.Id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	for _, track := range tracks {
		trackIds = append(trackIds, spotify.ID(track.Id))
	}
	if len(tracks) > 0 {
		_, err = client.AddTracksToPlaylist(playlistAPI.ID, trackIds...)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"href": playlistAPI.ExternalURLs["spotify"],
	})
}

func (h *Handler) AddTrackToSpotifyFav(ctx *gin.Context) {
	trackId := ctx.Param("id")
	token, err := GetSpotifyToken(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	client := authAPI.NewClient(token)
	userHas, err := client.UserHasTracks(spotify.ID(trackId))
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if userHas[0] {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("you already have this song in library").Error())
		return
	}
	err = client.AddTracksToLibrary(spotify.ID(trackId))
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) GetSpotifyProfileInfo(ctx *gin.Context) {
	var track music.Track
	var tracks []music.Track
	var artist music.Artist
	var artists []music.Artist
	token, err := GetSpotifyToken(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	client := authAPI.NewClient(token)

	artistsApi, err := client.CurrentUsersTopArtists()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	for _, artistAPI := range artistsApi.Artists {
		artist.NewArtist(artistAPI)
		artists = append(artists, artist)
		artist = music.Artist{}
	}
	tracksAPI, err := client.CurrentUsersTopTracks()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	for _, tr := range tracksAPI.Tracks {
		track.NewTrack(tr)
		tracks = append(tracks, track)
		track = music.Track{}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"artists": artists,
		"tracks":  tracks,
	})
}
