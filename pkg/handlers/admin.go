package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type JSONReason struct {
	Reason string `json:"reason" binding:"required"`
}

func (h *Handler) GetAdminUserPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin_user.html", nil)
}

func (h *Handler) GetAdminArtistPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin_artist.html", nil)
}

func (h *Handler) GetAdminPlaylistPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin_playlist.html", nil)
}

func (h *Handler) GetAdminReportPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin_report.html", nil)
}

func (h *Handler) GetAdminAnalyticsPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin_analytic.html", nil)
}

func (h *Handler) CheckAdmin(ctx *gin.Context) {
	userID, err := GetUserId(ctx)
	if err != nil {
		return
	}
	check, err := h.service.AdminService.CheckAdmin(userID)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if !check {
		newErrorResponse(ctx, http.StatusUnauthorized, "you are not admin")
		return
	}
}

func (h *Handler) GetAllUsers(ctx *gin.Context) {
	users, err := h.service.UserService.GetAllUsers()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func (h *Handler) GetAllPlaylistsForAdmin(ctx *gin.Context) {
	playlists, err := h.service.PlaylistService.GetAllPlaylistForAdmin()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"playlists": playlists,
	})
}

func (h *Handler) GetAllArtists(ctx *gin.Context) {
	artists, err := h.service.ArtistService.GetAllArtists()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"artists": artists,
	})
}

func (h *Handler) DeleteUser(ctx *gin.Context) {
	var input JSONReason
	userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	isAdmin, err := h.service.AdminService.CheckAdmin(userID)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if isAdmin {
		newErrorResponse(ctx, http.StatusBadRequest, errors.New("you can not delete admin").Error())
		return
	}
	err = h.service.UserService.DeleteUser(userID, input.Reason)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) DeleteArtist(ctx *gin.Context) {
	var input JSONReason
	artistID := ctx.Param("id")
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	err := h.service.ArtistService.DeleteArtist(artistID, input.Reason)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) DeletePlaylistByAdmin(ctx *gin.Context) {
	var input JSONReason
	playlistID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	err = h.service.PlaylistService.DeletePlaylistByAdmin(playlistID, input.Reason)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) GetHistory(ctx *gin.Context) {
	history, err := h.service.AdminService.GetHistory()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"history": history,
	})
}

func (h *Handler) GetReport(ctx *gin.Context) {
	var path string
	var err error
	format := ctx.Param("format")
	if format == "pdf" {
		path, err = h.service.AdminService.GetReportPDF()
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if format == "excel" {
		path, err = h.service.AdminService.GetReportExcel()
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}

	ctx.File(path)
	ctx.Status(http.StatusOK)
}

func (h *Handler) ExportJSON(ctx *gin.Context) {
	filePath, err := h.service.AdminService.GetDBInJSON()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.File(filePath)
	ctx.Status(http.StatusOK)
}

func (h *Handler) ImportJSON(ctx *gin.Context) {
	file, err := ctx.FormFile("json")
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.service.AdminService.ImportDBInJSON(file)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) GetGenrePopularity(ctx *gin.Context) {
	genres, err := h.service.AdminService.GetGenrePopularity()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"genres": genres,
	})
}

func (h *Handler) GetArtistPopularity(ctx *gin.Context) {
	artists, err := h.service.AdminService.GetArtistPopularity()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"artists": artists,
	})
}

func (h *Handler) GetGenreDiversity(ctx *gin.Context) {
	genres, err := h.service.AdminService.GetGenreDiversity()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"genres": genres,
	})
}
