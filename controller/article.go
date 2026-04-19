package controller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/fullstop113/go-web3-demo/model"
	"github.com/fullstop113/go-web3-demo/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	articleStatusDraft = "draft"
	articleStatusPublished = "published"
)

func CreateArticle(c *gin.Context) {
	var req CreateArticleRequest
	if err:= c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, utils.CodeInvalid, "invalid request body")
		return
	}
	req.Title = strings.TrimSpace(req.Title)
	req.Content = strings.TrimSpace(req.Content)
	req.Status = normalizeArticleStatus(req.Status)
	if req.Status == "" {
		req.Status = articleStatusDraft
	}
	if req.Title == "" || req.Content == "" || !isValidArticleStatus(req.Status) {
		utils.Fail(c, http.StatusBadRequest, utils.CodeInvalid, "invalid article fields")
	}

	userIDValue, ok := c.Get("user_id")
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, utils.CodeAuth, "unauthorized")
		return
	}
	userID, ok := userIDValue.(uint)
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, utils.CodeAuth, "unauthorized")
		return
	}
	article := model.Article{
		Title: req.Title,
		Content: req.Content,
		Status: req.Status,
		AuthorID: userID,
	}
	if err := model.DB.Create(&article).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, utils.CodeServer, "create article failed")
		return
	}
	utils.OK(c, article)
}

func ListArticles(c *gin.Context) {
	page := parsePositiveInt(c.DefaultQuery("page", "1"), 1)
	pageSize := parsePositiveInt(c.DefaultQuery("page_size", "10"), 10)
	if pageSize > 100 {
		pageSize = 100
	}
	var total int64
	if err := model.DB.Model(&model.Article{}).Count(&total).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, utils.CodeServer, "query articles failed")
		return
	}
	var articles []model.Article
	offset := (page - 1) * pageSize
	if err := model.DB.Order("id desc").Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, utils.CodeInvalid, "query articles failed")
		return
	}
	utils.OK(c, gin.H{
		"list": articles,
		"page": page,
		"page_size": pageSize,
		"total": total,
	})
}

func GetArticle(c *gin.Context) {
	var article model.Article
	if err := model.DB.First(&article, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Fail(c, http.StatusNotFound, utils.CodeNotFound, "article not found")
			return
		}
		utils.Fail(c, http.StatusInternalServerError, utils.CodeServer, "query article failed")
		return
	}
	utils.OK(c, article)
}

func UpdateArticle(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, utils.CodeAuth, "unauthorized")
		return
	}
	var article model.Article
	if err := model.DB.First(&article, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Fail(c, http.StatusNotFound, utils.CodeNotFound, "article not found")
			return
		}
		utils.Fail(c, http.StatusInternalServerError, utils.CodeServer, "query article failed")
		return
	}
	if article.AuthorID != userID {
		utils.Fail(c, http.StatusForbidden, utils.CodeForbidden, "forbidden")
		return
	}
	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, utils.CodeInvalid, "invaild requrest body")
		return
	}
	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title == "" {
			utils.Fail(c, http.StatusBadRequest, utils.CodeInvalid, "title cannot be empty")
			return
		}
		article.Title = title
	}
	if req.Content != nil {
		content := strings.TrimSpace(*req.Content)
		if content == "" {
			utils.Fail(http.StatusBadRequest, utils.CodeInvalid, "content not be empty")
			return
		}
		article.Content = content
	}
	if req.Status != nil {
		status := normalizeArticleStatus(*req.Status)
		if !isValidArticleStatus(status) {
			utils.Fail(c, http.StatusBadRequest, utils.CodeInvalid, "invalid status")
			return
		}
		article.Status = status
	}
	if err := model.DB.Save(&article).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, utils.CodeServer, "update article failed")
		return
	}
	utils.OK(c, article)
}

func DeleteArticle(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, utils.CodeAuth, "unauthorized")
		return
	}

	var article model.Article
	if err := model.DB.First(&article, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Fail(c, http.StatusNotFound, utils.CodeNotFound, "article not found")
			return
		}
		utils.Fail(c, http.StatusInternalServerError, utils.CodeServer, "query article failed")
		return
	}
	if article.AuthorID != userID {
		utils.Fail(c, http.StatusForbidden, utils.CodeForbidden, "forbidden")
		return
	}
	if err := model.DB.Delete(&article).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, utils.CodeServer, "delete article failed")
		return
	}
	utils.OK(c, gin.H{"deleted": true})
}

func getUserIDFromContext(c *gin.Context) (uint, bool) {
	value, ok := c.Get("user_id")
	if !ok {
		return 0, false
	}
	userID, ok := value.(uint)
	if !ok {
		return 0, false
	}
	return userID, true
}

func isValidArticleStatus(status string) bool {
	return status == articleStatusDraft || status == articleStatusPublished
}

func normalizeArticleStatus(status string) string {
	return strings.ToLower(strings.TrimSpace(status))
}

func parsePositiveInt(raw string, defaultValue int) int {
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return defaultValue
	}
	return value
}
