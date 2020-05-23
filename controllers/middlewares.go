package controllers

import (
	"net/http"

	"github.com/bin16/tudou/db"
	"github.com/bin16/tudou/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// ShouldAuthed pick user from token
func ShouldAuthed(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	tokenString := header[7:]
	claims := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if !token.Valid {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	email := claims.Id

	u0 := models.User{Email: email}
	if db.DB.Preload("Setting").First(&u0).Error != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	c.Set("user", &u0)
	c.Next()
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	tokenRecord := models.TokenRecord{Token: token.Raw}
	if err := db.DB.First(&tokenRecord).Error; err != nil {
		return []byte{}, err
	}

	return tokenRecord.Key, nil
}
