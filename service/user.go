package service

import (
	"errors"

	"github.com/fullstop113/go-web3-demo/model"
	"github.com/fullstop113/go-web3-demo/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct{}

// Register 用户注册
func (s *UserService) Register(username, password, email string) (*model.User, error) {
	// 检查用户名是否已存在
	var existUser model.User
	if err := model.DB.Where("username = ?", username).First(&existUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if err := model.DB.Where("email = ?", email).First(&existUser).Error; err == nil {
		return nil, errors.New("邮箱已被注册")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := model.DB.Create(user).Error; err != nil {
		return nil, errors.New("创建用户失败")
	}

	return user, nil
}

// Login 用户登录，返回token
func (s *UserService) Login(username, password string) (string, error) {
	var user model.User
	if err := model.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("用户不存在")
		}
		return "", errors.New("数据库查询失败")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("密码错误")
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", errors.New("生成token失败")
	}

	return token, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := model.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
