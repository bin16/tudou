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
		c.JSON(msgJSON(http.StatusForbidden, "token_not_found"))
		c.Abort()
		return
	}

	tokenString := header[7:]
	claims := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc)
	if err != nil {
		c.JSON(msgJSON(http.StatusBadRequest, "bad_token"))
		c.Abort()
		return
	}

	if !token.Valid {
		c.JSON(msgJSON(http.StatusForbidden, "invalid_token"))
		c.Abort()
		return
	}

	u0 := models.User{}
	if db.DB.Preload("Setting").Where("uuid = ?", claims.Id).First(&u0).Error != nil {
		c.JSON(msgJSON(http.StatusForbidden, "user_not_found"))
		c.Abort()
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
