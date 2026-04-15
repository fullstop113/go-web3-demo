package controller

import (
	"net/http"

	"github.com/fullstop113/go-web3-demo/utils"
	"github.com/gin-gonic/gin"
)


func GetUserInfo(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, utils.CodeAuth, "unauthorized")
		return
	}

	username, _ := c.Get("username")
	utils.OK(c, gin.H{
		"user_id": userID,
		"username": username,
	})
}