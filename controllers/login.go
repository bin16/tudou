package controllers

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/bin16/tudou/db"
	"github.com/bin16/tudou/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// PostLogin handle login form
func PostLogin(c *gin.Context) {
	email := c.PostForm("email")
	user := models.User{Email: email}
	if db.DB.First(&user).Error != nil {
		c.Status(http.StatusForbidden)
		return
	}

	key := make([]byte, 64)
	rand.Read(key)

	claims := &jwt.StandardClaims{
		Audience:  "tudou",
		Issuer:    "tudou",
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().AddDate(0, 1, 0).Unix(),
		Subject:   "login",
		Id:        email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenS, err := token.SignedString(key)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}

	if err := db.DB.Create(&models.TokenRecord{
		Token:   tokenS,
		Key:     key,
		Revoked: false,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenS,
	})
}
