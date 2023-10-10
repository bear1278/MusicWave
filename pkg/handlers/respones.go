package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
)

type error struct {
	Message string `json:"message"`
}

func newErrorResponse(ctx *gin.Context, statuscode int, message string) {
	log.Printf(message)
	ctx.AbortWithStatusJSON(statuscode, error{message})
}
