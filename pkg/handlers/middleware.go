package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader        = "Authorization-1"
	authorizationSpotifyHeader = "Authorization-2"
	refresh                    = "Authorization-3"
	expiry                     = "Authorization-4"
	userCtx                    = "userId"
	spotifyCtx                 = "accessToken"
	refreshCtx                 = "refreshToken"
	expiryCtx                  = "expiry"
)

func (h *Handler) userIdentity(ctx *gin.Context) {
	header := ctx.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(ctx, http.StatusUnauthorized, "empty auth header")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(ctx, http.StatusUnauthorized, "invalid auth header")
		return
	}

	UserId, err := h.service.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
	}

	ctx.Set(userCtx, UserId)
}

func (h *Handler) spotifyIdentity(ctx *gin.Context) {
	header := ctx.GetHeader(authorizationSpotifyHeader)
	header1 := ctx.GetHeader(refresh)
	header2 := ctx.GetHeader(expiry)
	if header == "" || header1 == "" || header2 == "" {
		newErrorResponse(ctx, http.StatusUnauthorized, "empty auth header")
		return
	}
	headerParts := strings.Split(header, " ")
	headerParts1 := strings.Split(header1, " ")
	headerParts2 := strings.Split(header2, " ")
	if len(headerParts) != 2 || len(headerParts1) != 2 || len(headerParts2) != 2 {
		newErrorResponse(ctx, http.StatusUnauthorized, "invalid auth header")
		return
	}

	ctx.Set(spotifyCtx, headerParts[1])
	ctx.Set(refreshCtx, headerParts1[1])
	ctx.Set(expiryCtx, headerParts2[1])
}
