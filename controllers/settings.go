package controllers

import (
	"net/http"

	"github.com/bin16/tudou/db"
	"github.com/bin16/tudou/models"
	"github.com/gin-gonic/gin"
)

// GetSettings GET /api/my/settings
func GetSettings(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	settings := models.Setting{UserID: user.ID}
	if err := db.DB.FirstOrCreate(&settings).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":     user,
		"settings": settings,
	})
}

// UpdateSettings POST /api/a/settings.update
func UpdateSettings(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	settings := models.Setting{UserID: user.ID}
	if err := c.BindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	settings.UserID = user.ID
	if err := db.DB.Model(&settings).Updates(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}
