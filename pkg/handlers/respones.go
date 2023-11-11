package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
