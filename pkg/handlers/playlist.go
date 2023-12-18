package handlers

import (
	"errors"
	"fmt"
	music "github.com/bear1278/MusicWave"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PlaylistJSON struct {
	Name  string `json:"name" binding:"required"`
	Cover string `json:"cover"`
	Type  string `json:"type" binding:"required"`
}

func (h *Handler) GetPlaylistPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "playlist.html", nil)
}

func (h *Handler) NewPlaylist(ctx *gin.Context) {
	var pathImages string
	var playlist music.Playlist
	userId, err := GetUserId(ctx)
	if err != nil {
		return
	}
	file, err := ctx.FormFile("cover")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			pathImages = "/static/images/music.svg"
		} else {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if err == nil {
		path := fmt.Sprintf("./web/images/playlists/%s", file.Filename)
		err = ctx.SaveUploadedFile(file, path)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		pathImages = fmt.Sprintf("/static/images/playlists/%s", file.Filename)
	}
	playlist.Name = ctx.PostForm("name")
	playlist.Cover = pathImages
	playlist.Type = ctx.PostForm("type")
	playlistId, err := h.service.PlaylistService.NewPlaylist(playlist, userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": playlistId,
	})
}

func (h *Handler) DeletePlaylist(ctx *gin.Context) {
	playlistID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	userID, err := GetUserId(ctx)
	if err != nil {
		return
	}
	if err := h.service.PlaylistService.DeletePlaylist(playlistID, userID); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

type getAllPlaylistsResponse struct {
	Favorites      []music.Playlist `json:"favorites"`
	MyPlaylists    []music.Playlist `json:"myPlaylists"`
	AddedPlaylists []music.Playlist `json:"addedPlaylists"`
}

func (h *Handler) GetAllPlaylistsForUser(ctx *gin.Context) {
	var favoritesSlice []music.Playlist
	userId, err := GetUserId(ctx)
	if err != nil {
		return
	}
	myPlaylists, addedPlaylists, err := h.service.PlaylistService.GetAllPlaylist(userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	favorites, err := h.service.PlaylistService.GetUserFavorites(userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	favoritesSlice = append(favoritesSlice, favorites)
	ctx.JSON(http.StatusOK, getAllPlaylistsResponse{
		Favorites:      favoritesSlice,
		MyPlaylists:    myPlaylists,
		AddedPlaylists: addedPlaylists,
	})
}

func (h *Handler) GetById(ctx *gin.Context) {
	var changePermission, addSpotifyPermission bool
	userId, err := GetUserId(ctx)
	if err != nil {

		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	playlist, err := h.service.PlaylistService.GetById(id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if playlist.Name != "favorites" && playlist.Author.Id == userId {
		changePermission = true
	} else {
		changePermission = false
	}

	if playlist.Name == "favorites" && playlist.Author.Id == userId {
		addSpotifyPermission = true
	} else {
		addSpotifyPermission = false
	}

	ctx.JSON(http.StatusOK, gin.H{
		"playlist":             playlist,
		"changePermission":     changePermission,
		"addSpotifyPermission": addSpotifyPermission,
	})
}

type PlaylistUpdateJSON struct {
	Name  string `json:"name"`
	Cover string `json:"cover"`
	Type  string `json:"type"`
}

func (h *Handler) UpdatePlaylist(ctx *gin.Context) {
	var playlist music.Playlist
	var pathImages string
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	oldPlaylist, err := h.service.PlaylistService.GetById(id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	file, err := ctx.FormFile("cover")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			pathImages = oldPlaylist.Cover
		} else {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if err == nil {
		path := fmt.Sprintf("./web/images/playlists/%s", file.Filename)
		err = ctx.SaveUploadedFile(file, path)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		pathImages = fmt.Sprintf("/static/images/playlists/%s", file.Filename)
	}
	playlist.Name = ctx.PostForm("name")
	playlist.Cover = pathImages
	playlist.Type = ctx.PostForm("type")

	if playlist.Name == "" {
		playlist.Name = oldPlaylist.Name
	}
	if playlist.Cover == "" {
		playlist.Cover = oldPlaylist.Cover
	}
	if playlist.Type == "" {
		playlist.Type = oldPlaylist.Type
	}
	playlist.Id = id
	playlist.Author.Id, err = GetUserId(ctx)
	if err != nil {
		return
	}
	if err := h.service.PlaylistService.UpdatePlaylist(playlist); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) ExcludePlaylist(ctx *gin.Context) {
	idPlaylist, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	idUser, err := GetUserId(ctx)
	if err != nil {
		return
	}
	err = h.service.PlaylistService.ExcludePlaylist(idPlaylist, idUser)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) AddPlaylist(ctx *gin.Context) {
	idPlaylist, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	idUser, err := GetUserId(ctx)
	if err != nil {
		return
	}
	err = h.service.PlaylistService.AddPlaylist(idPlaylist, idUser)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) SearchPlaylist(ctx *gin.Context) {
	searchString := ctx.Param("string")
	playlists, err := h.service.PlaylistService.SearchPlaylist(searchString)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": playlists})
}

func (h *Handler) GetUsersPlaylists(ctx *gin.Context) {
	userId, err := GetUserId(ctx)
	if err != nil {
		return
	}
	myPlaylists, _, err := h.service.PlaylistService.GetAllPlaylist(userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	favorites, err := h.service.PlaylistService.GetUserFavorites(userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	myPlaylists = append(myPlaylists, favorites)
	ctx.JSON(http.StatusOK, gin.H{
		"playlists": myPlaylists,
	})
}
