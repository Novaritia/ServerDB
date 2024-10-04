package handler

import (
	"context"
	"net/http"
	"serverdb/ent"

	"github.com/gin-gonic/gin"
)

type CreateUserStruct struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var client *ent.Client

func RegisterRoute(router *gin.Engine, dbClient *ent.Client) {
	client = dbClient
	router.GET("/ping", Ping)
	router.POST("/users", CreateUser)
	router.NoRoute(HandleError)
}

func Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func CreateUser(ctx *gin.Context) {
	var req CreateUserStruct

	if ctx.Request.ContentLength == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request body is empty"})
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser, err := client.User.Create().
		SetName(req.Name).
		SetPassword(req.Password).
		Save(context.Background())

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user!"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":       newUser.ID,
		"name":     newUser.Name,
		"password": newUser.Password,
	})
}

func HandleError(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"code":    404,
		"message": "incorrect path data not found!",
	})
}
