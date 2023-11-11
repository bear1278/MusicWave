package handlers

import (
	music "github.com/bear1278/MusicWave"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PlaylistJSON struct {
	Name  string `json:"name" binding:"required"`
	Cover string `json:"cover" binding:"required"`
}

func (h *Handler) NewPlaylist(ctx *gin.Context) {
	var input PlaylistJSON
	var playlist music.Playlist
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	userId, err := GetUserId(ctx)
	if err != nil {
		return
	}
	playlist.Name = input.Name
	playlist.Cover = input.Cover
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
	var playlist music.Playlist
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	playlist.Id = id
	if err := h.service.PlaylistService.DeletePlaylist(playlist.Id); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

type getAllPlaylistsResponse struct {
	Data []music.Playlist `json:"data"`
}

func (h *Handler) GetAllPlaylists(ctx *gin.Context) {
	userId, err := GetUserId(ctx)
	if err != nil {
		return
	}
	playlists, err := h.service.PlaylistService.GetAllPlaylist(userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, getAllPlaylistsResponse{
		Data: playlists,
	})
}

func (h *Handler) GetById(ctx *gin.Context) {
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
	ctx.JSON(http.StatusOK, playlist)
}

func (h *Handler) UpdatePlaylist(ctx *gin.Context) {
	var input PlaylistJSON
	var playlist music.Playlist
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	playlist.Id = id
	playlist.Name = input.Name
	playlist.Cover = input.Cover
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
