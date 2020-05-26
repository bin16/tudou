package token

import (
	"crypto/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/bin16/tudou/db"
	"github.com/bin16/tudou/models"
	"github.com/google/uuid"

	"github.com/dgrijalva/jwt-go"
)

func CheckAccessToken(c *gin.Context) {
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
	if db.DB.Preload("Setting").Where("uuid = ?", claims.Subject).First(&u0).Error != nil {
		c.JSON(msgJSON(http.StatusForbidden, "user_not_found"))
		c.Abort()
		return
	}

	c.Set("user", &u0)
	c.Next()
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	claims := token.Claims.(*jwt.StandardClaims)

	tokenRecord := models.TokenRecord{}
	if err := db.DB.Where("jwt_id = ?", claims.Id).First(&tokenRecord).Error; err != nil {
		return []byte{}, err
	}

	return tokenRecord.Key, nil
}

func msgJSON(status int, message string) (int, gin.H) {
	return status, gin.H{"message": message}
}

func AccessToken(subject string) (string, error) {
	issuedAt := time.Now().Unix()
	expiresAt := time.Now().AddDate(0, 1, 0).Unix()
	issuer := "tudou"
	audience := "tudou"
	jwtID := uuid.New()
	claims := jwt.StandardClaims{
		Audience:  audience,
		Issuer:    issuer,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		Subject:   subject,
		Id:        jwtID.String(),
	}

	key := make([]byte, 64)
	rand.Read(key)

	err := db.DB.Create(&models.TokenRecord{
		Key:       key,
		JwtID:     jwtID.String(),
		Revoked:   false,
		Subject:   subject,
		ExpiresAt: expiresAt,
		IssuedAt:  issuedAt,
		Issuer:    issuer,
		Audience:  audience,
	}).Error
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}
