package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"time"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(ctx *gin.Context, statuscode int, message string) {
	log.Printf(message)
	ctx.AbortWithStatusJSON(statuscode, errorResponse{message})
}

func GetUserId(ctx *gin.Context) (int64, error) {
	userId, ok := ctx.Get(userCtx)
	if !ok {
		newErrorResponse(ctx, http.StatusInternalServerError, "id not found")
		return 0, errors.New("id not found")
	}
	id, ok := userId.(int64)
	if !ok {
		newErrorResponse(ctx, http.StatusInternalServerError, "id is invalid type")
		return 0, errors.New("id is invalid type")
	}
	return id, nil
}

func GetSpotifyToken(ctx *gin.Context) (*oauth2.Token, error) {
	tokenSpotify, ok := ctx.Get(spotifyCtx)
	if !ok {
		newErrorResponse(ctx, http.StatusInternalServerError, "id not found")
		return nil, errors.New("id not found")
	}
	refresh, ok := ctx.Get(refreshCtx)
	if !ok {
		newErrorResponse(ctx, http.StatusInternalServerError, "id not found")
		return nil, errors.New("id not found")
	}
	expiry, ok := ctx.Get(expiryCtx)
	if !ok {
		newErrorResponse(ctx, http.StatusInternalServerError, "id not found")
		return nil, errors.New("id not found")
	}
	accessToken, ok := tokenSpotify.(string)
	if !ok {
		newErrorResponse(ctx, http.StatusInternalServerError, "id is invalid type")
		return nil, errors.New("id is invalid type")
	}
	refreshToken, ok := refresh.(string)
	if !ok {
		newErrorResponse(ctx, http.StatusInternalServerError, "id is invalid type")
		return nil, errors.New("id is invalid type")
	}
	expiryTime, ok := expiry.(string)
	if !ok {
		newErrorResponse(ctx, http.StatusInternalServerError, "id is invalid type")
		return nil, errors.New("id is invalid type")
	}
	duration, err := time.Parse(time.RFC3339, expiryTime)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "id is invalid type")
		return nil, errors.New("id is invalid type")
	}

	return &oauth2.Token{AccessToken: accessToken, RefreshToken: refreshToken, Expiry: duration}, nil
}
