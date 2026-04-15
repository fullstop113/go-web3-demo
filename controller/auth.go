package controller

import (
	"net/http"
	"time"

	"github.com/fullstop113/go-web3-demo/model"
	"github.com/fullstop113/go-web3-demo/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, utils.CodeAuth, "invalid request body")
		return
	}

	// TODO: 用户名/邮箱查重
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, utils.CodeServer, "hash password failed")
		return
	}
	u := model.User{
		Username: req.Username,
		Email:    req.Email,
		PasswordHash: string(hash),
	}
	if err := model.DB.Create(&u).Error; err != nil {
		utils.Fail(c, http.StatusBadRequest, utils.CodeInvalid, "create user failed")
		return
	}
	utils.OK(c, gin.H{
		"id":       u.ID,
		"username": u.Username,
		"email": u.Email,
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, utils.CodeInvalid, "invaild requrest")
		return
	}
	var u model.User
	if err := model.DB.Where("username = ?", req.Username).First(&u).Error; err != nil {
		utils.Fail(c, http.StatusUnauthorized, utils.CodeAuth, "invalid username or password")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password),); err != nil {
		utils.Fail(c, http.StatusUnauthorized, utils.CodeAuth, "invalid username or password")
		return
	}

	// 需要我在utils/jwt.go 增加GenerateToken(userID, username)
	token, expireAt, err := utils.GenerateToken(u.ID, u.Username)
	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, utils.CodeServer, "generate token failed")
		return
	}
	utils.OK(c, gin.H{
		"token_type": "Bearer",
		"access_token": token,
		"expires_at": expireAt.Format(time.RFC3339),
	})

}