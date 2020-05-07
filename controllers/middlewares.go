package controllers

import (
	"net/http"

	"github.com/bin16/tudou/db"
	"github.com/bin16/tudou/models"
	"github.com/gin-gonic/gin"
)

// ShouldAuthed pick user from token
func ShouldAuthed(c *gin.Context) {
	u0 := models.User{}
	if db.DB.Debug().Preload("Setting").First(&u0).Error != nil {
		c.Status(http.StatusForbidden)
		return
	}

	c.Set("user", &u0)
	c.Next()
}
