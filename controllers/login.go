package controllers

import (
	"net/http"
	"time"

	"github.com/bin16/tudou/db"
	"github.com/bin16/tudou/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func PostLogin(c *gin.Context) {
	email := c.PostForm("email")
	user := models.User{Email: email}
	if db.DB.First(&user).Error != nil {
		c.Status(http.StatusForbidden)
		return
	}

	claims := &jwt.StandardClaims{
		Audience:  "tudou",
		Issuer:    "tudou",
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().AddDate(0, 1, 0).Unix(),
		Subject:   "login",
		Id:        email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenS, err := token.SignedString([]byte("h908H*Hjhoiwh809h78gs6^$E^TRCVJM(82y,l;L?<O@E)Y0"))
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenS,
	})
}
