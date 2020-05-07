package controllers

import (
	"net/http"

	"github.com/bin16/tudou/controllers/forms"
	"github.com/bin16/tudou/db"
	"github.com/bin16/tudou/models"
	"github.com/gin-gonic/gin"
)

// CreateUser POST /a/user.create
func CreateUser(c *gin.Context) {
	values := forms.SignupForm{}
	if err := c.BindJSON(&values); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if valid, message := values.Valid(); !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	u := values.ToUser()
	if err := db.DB.Create(&u).Error; err != nil {
		c.JSON(errorJSON(err))
		return
	}

	c.JSON(errorJSON(nil))
}

func GetUsers(c *gin.Context) {
	ul := []models.User{}
	if err := db.DB.Find(&ul).Error; err != nil {
		c.JSON(errorJSON(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": ul,
	})
}
