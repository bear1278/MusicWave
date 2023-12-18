package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *Handler) GetMainPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "main.html", nil)
}
func (h *Handler) GetRelocatePage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "relocate.html", nil)
}

func (h *Handler) GetLibraryPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "library.html", nil)
}

func (h *Handler) GetUserRecommendation(ctx *gin.Context) {
	userId, err := GetUserId(ctx)
	if err != nil {
		return
	}
	tracks, err := h.service.TrackService.GetTracksForRec(userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	albums, err := h.service.AlbumService.GetAlbumsRec(userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	for k, _ := range albums {
		if albums[k].Duration == 0 {
			duration, err := h.GetAlbumDurationFromApi(albums[k].Id, ctx)
			if err != nil {
				newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
				return
			}
			albums[k].Duration = duration
		}
	}
	artists, err := h.service.ArtistService.GetArtistsRec(userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	playlists, err := h.service.PlaylistService.GerPlaylistsRec(userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"tracks":    tracks,
		"albums":    albums,
		"artists":   artists,
		"playlists": playlists,
	})
	log.Printf("%s", ctx.Errors)
}
