package controller

import (
	"net/http"

	"github.com/fullstop113/go-web3-demo/middleware"
	"github.com/fullstop113/go-web3-demo/service"
	"github.com/fullstop113/go-web3-demo/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: &service.UserService{},
	}
}

// Register 用户注册
func (ctrl *UserController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, utils.CodeBadRequest, "请求参数错误")
		return
	}

	user, err := ctrl.userService.Register(req.Username, req.Password, req.Email)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, utils.CodeBadRequest, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

// Login 用户登录
func (ctrl *UserController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, utils.CodeBadRequest, "请求参数错误")
		return
	}

	token, err := ctrl.userService.Login(req.Username, req.Password)
	if err != nil {
		utils.Fail(c, http.StatusUnauthorized, utils.CodeUnauthorized, err.Error())
		return
	}

	utils.OK(c, gin.H{
		"token": token,
	})
}

// GetCurrentUser 获取当前用户信息
func (ctrl *UserController) GetCurrentUser(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)

	user, err := ctrl.userService.GetUserByID(userID)
	if err != nil {
		utils.Fail(c, http.StatusNotFound, utils.CodeNotFound, "用户不存在")
		return
	}

	utils.OK(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}
