package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type JSONUsername struct {
	Username string `json:"username" binding:"required"`
}

type JSONPassword struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

type JSONPicture struct {
	Picture string `json:"picture" binding:"required"`
}

type JSONEmail struct {
	Email string `json:"email" binding:"required"`
}

func (h *Handler) ChangeUsername(ctx *gin.Context) {
	var input JSONUsername
	userID, err := GetUserId(ctx)
	if err != nil {
		return
	}
	if err = ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err = h.service.UserService.ChangeUsername(userID, input.Username); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) ChangePassword(ctx *gin.Context) {
	var input JSONPassword

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	userID, err := GetUserId(ctx)
	if err != nil {
		return
	}
	if err = h.service.UserService.ChangePassword(userID, h.service.Authorization.GeneratePasswordHash(input.NewPassword),
		h.service.Authorization.GeneratePasswordHash(input.OldPassword)); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) ChangePicture(ctx *gin.Context) {
	userID, err := GetUserId(ctx)
	if err != nil {
		return
	}
	file, err := ctx.FormFile("picture")
	if err != nil {

		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return

	}
	path := fmt.Sprintf("./web/images/users/%s", file.Filename)
	err = ctx.SaveUploadedFile(file, path)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	pathImages := fmt.Sprintf("/static/images/users/%s", file.Filename)

	if err = h.service.UserService.ChangePicture(userID, pathImages); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) ChangeEmail(ctx *gin.Context) {
	var input JSONEmail
	userID, err := GetUserId(ctx)
	if err != nil {
		return
	}
	if err = ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err = h.service.UserService.ChangeEmail(userID, input.Email); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *Handler) GetProfilePage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "profile.html", nil)
}

func (h *Handler) GetUserForProfile(ctx *gin.Context) {
	userId, err := GetUserId(ctx)
	if err != nil {
		return
	}
	user, err := h.service.UserService.GetUserByID(userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (h *Handler) GetProfileInfo(ctx *gin.Context) {
	userId, err := GetUserId(ctx)
	if err != nil {
		return
	}
	tracks, err := h.service.TrackService.GetMostPopularTracksForUser(userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	artists, err := h.service.ArtistService.GetMostPopularTrackForUser(userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	genres, err := h.service.GenreService.GerUserGenres(userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	duration, err := h.service.PlaylistService.GetDurationOfAllPlaylists(userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"tracks":   tracks,
		"artists":  artists,
		"genres":   genres,
		"duration": duration,
	})
}
