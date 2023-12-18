package handlers

import (
	music "github.com/bear1278/MusicWave"
	"github.com/gin-gonic/gin"
	"net/http"
)

type JSONArtistResponse struct {
	Data []music.Artist `json:"artists"`
}

func (h *Handler) GetArtistPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "artist.html", nil)
}

func (h *Handler) GetAllArtistForUser(ctx *gin.Context) {
	idUser, err := GetUserId(ctx)
	if err != nil {
		return
	}
	artist, err := h.service.ArtistService.GetAllAlbumsForUser(idUser)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, JSONArtistResponse{
		Data: artist,
	})
}

func (h *Handler) GetArtistById(ctx *gin.Context) {
	idArtist := ctx.Param("id")
	artist, err := h.service.ArtistService.GetArtistById(idArtist)
	if artist.Id == "" {
		artist, err = h.GetArtistByIDFromApi(idArtist, ctx)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"artist": artist,
	})
}

func (h *Handler) ExcludeArtist(ctx *gin.Context) {
	idUser, err := GetUserId(ctx)
	if err != nil {
		return
	}
	idArtist := ctx.Param("id")
	if err = h.service.ArtistService.ExcludeArtist(idUser, idArtist); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) AddArtistToFav(ctx *gin.Context) {
	var artist music.Artist
	if err := ctx.BindJSON(&artist); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	idUser, err := GetUserId(ctx)
	if err != nil {
		return
	}
	err = h.service.ArtistService.AddArtistToFav(idUser, artist)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	tracks, err := h.GetTop5TracksOfArtist(artist.Id, ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	for _, tr := range tracks {
		album, err := h.GetAlbumByIDFromApi(tr.Album.Id, ctx)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		if len(album.Artist) > 1 {
			for _, secondArtist := range album.Artist {
				if secondArtist.Id != artist.Id {
					secondArtist, err = h.GetArtistByIDFromApi(secondArtist.Id, ctx)
					if err != nil {
						newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
						return
					}
					err = h.service.AddArtistToDB(secondArtist)
					if err != nil {
						newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
						return
					}
				}
			}
		}
		err = h.service.AlbumService.AddAlbumToDB(album)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		_, err = h.service.TrackService.AddTrackToDB(tr)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}
	err = h.service.GenreService.AddGenreToArtist(artist.Id, artist.Genres)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}
