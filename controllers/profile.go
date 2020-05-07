package controllers

import (
	"net/http"

	"github.com/bin16/tudou/models"
	"github.com/gin-gonic/gin"
)

// GetProfile GET /api/my/profile
func GetProfile(c *gin.Context) {
	u0 := c.MustGet("user").(*models.User)
	c.JSON(http.StatusOK, gin.H{
		"user": u0,
	})
}
