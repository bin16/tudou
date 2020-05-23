package controllers

import (
	"net/http"
	"time"

	"github.com/bin16/tudou/db"
	"github.com/bin16/tudou/models"
	"github.com/gin-gonic/gin"
)

// GetLogs ...
func GetLogs(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	logs := []models.Log{}
	yestaday := time.Now().AddDate(0, 0, -1)
	db.DB.Preload("Event").Where("user_id = ?", user.ID).Where("created_at > ?", yestaday).Order("id DESC").Find(&logs)
	c.JSON(http.StatusOK, gin.H{
		"coins": user.Coins,
		"logs":  logs,
		"user":  *user,
	})
}
