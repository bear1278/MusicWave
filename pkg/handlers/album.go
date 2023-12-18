package handlers

import (
	music "github.com/bear1278/MusicWave"
	"github.com/gin-gonic/gin"
	"net/http"
)

type JSONAlbumResponse struct {
	Data []music.Album `json:"albums"`
}

func (h *Handler) GetAlbumPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "album.html", nil)
}

func (h *Handler) GetAlbumById(ctx *gin.Context) {
	var album music.Album
	idAlbum := ctx.Param("id")

	album, err := h.service.AlbumService.GetAlbumById(idAlbum)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if album.Id == "" {
		album, err = h.GetAlbumByIDFromApi(idAlbum, ctx)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}
	album.Duration, err = h.GetAlbumDurationFromApi(album.Id, ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"album": album,
	})
}

func (h *Handler) GetAllAlbumsForUser(ctx *gin.Context) {
	idUser, err := GetUserId(ctx)
	if err != nil {
		return
	}
	albums, err := h.service.AlbumService.GetAllAlbumsForUser(idUser)
	for k, _ := range albums {
		albums[k].Duration, err = h.GetAlbumDurationFromApi(albums[k].Id, ctx)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, JSONAlbumResponse{
		Data: albums,
	})
}

func (h *Handler) GetAlbumsByArtist(ctx *gin.Context) {
	idArtist := ctx.Param("id")
	albums, err := h.GetAlbumsOfArtistFromApi(idArtist, ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	for k, _ := range albums {
		albums[k].Duration, err = h.GetAlbumDurationFromApi(albums[k].Id, ctx)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}
	ctx.JSON(http.StatusOK, JSONAlbumResponse{
		Data: albums,
	})
}

func (h *Handler) ExcludeAlbum(ctx *gin.Context) {
	idUser, err := GetUserId(ctx)
	if err != nil {
		return
	}
	idAlbum := ctx.Param("id")

	if err = h.service.AlbumService.ExcludeAlbum(idUser, idAlbum); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) AddAlbumToFav(ctx *gin.Context) {
	var album music.Album
	if err := ctx.BindJSON(&album); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	idUser, err := GetUserId(ctx)
	if err != nil {
		return
	}
	for _, artist := range album.Artist {
		artistAPI, err := h.GetArtistByIDFromApi(artist.Id, ctx)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		err = h.service.ArtistService.AddArtistToDB(artistAPI)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}
	albumAPI, err := h.GetAlbumByIDFromApi(album.Id, ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	err = h.service.AlbumService.AddAlbumToFav(idUser, albumAPI)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	tracks, err := h.GetTracksFromAlbumFromApi(album.Id, ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	for _, track := range tracks {
		_, err = h.service.TrackService.AddTrackToDB(track)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}
	ctx.Status(http.StatusOK)
}
