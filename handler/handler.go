package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoute(router *gin.Engine) {
	router.GET("/ping", Ping)
	router.NoRoute(HandleError)
}

func Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func HandleError(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    404,
		"message": "incorrect path data not found!",
	})
}
