package handlers

import (
	music "github.com/bear1278/MusicWave"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type JSONTracksResponse struct {
	Data []music.Track `json:"tracks"`
}

func (h *Handler) GetTrackPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "track.html", nil)
}

func (h *Handler) AddTrack(ctx *gin.Context) {
	var track music.Track
	if err := ctx.BindJSON(&track); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	idPlaylist, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	for _, artist := range track.Album.Artist {
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
		err = h.service.GenreService.AddGenreToArtist(artistAPI.Id, artistAPI.Genres)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}
	albumAPI, err := h.GetAlbumByIDFromApi(track.Album.Id, ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	err = h.service.AlbumService.AddAlbumToDB(albumAPI)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	_, err = h.service.TrackService.AddTrack(idPlaylist, track)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) ExcludeTrack(ctx *gin.Context) {
	idPlaylist, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	idTrack := ctx.Param("id_track")
	if err := h.service.TrackService.ExcludeTrack(idPlaylist, idTrack); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) GetAllTracks(ctx *gin.Context) {
	idPlaylist, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	tracks, err := h.service.TrackService.GetAllTracks(idPlaylist)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, JSONTracksResponse{
		Data: tracks,
	})
}

func (h *Handler) GetTrackById(ctx *gin.Context) {
	idTrack := ctx.Param("id_track")
	track, err := h.service.TrackService.GetTrackById(idTrack)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if track.Id == "" {
		track, err = h.GetTrackByIDFromAPI(idTrack, ctx)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"track": track,
	})
}

func (h *Handler) GetTracksFromAlbum(ctx *gin.Context) {
	idAlbum := ctx.Param("id")

	tracks, err := h.service.TrackService.GetTrackFromAlbum(idAlbum)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
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
	if int64(len(tracks)) < album.TotalTracks {
		tracks, err = h.GetTracksFromAlbumFromApi(idAlbum, ctx)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}
	ctx.JSON(http.StatusOK, JSONTracksResponse{
		Data: tracks,
	})
}

func (h *Handler) GetTopTracksByArtist(ctx *gin.Context) {
	idArtist := ctx.Param("id")
	tracks, err := h.service.TrackService.GetTopTracksOfArtist(idArtist)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if len(tracks) < 4 {
		tracks, err = h.GetTop5TracksOfArtist(idArtist, ctx)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"tracks": tracks,
	})
}
