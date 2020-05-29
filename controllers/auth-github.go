package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bin16/tudou/conf"
	"github.com/bin16/tudou/db"
	"github.com/bin16/tudou/models"
	"github.com/bin16/tudou/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type githubProfile struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

// OAuth2Github Github OAuth2
func OAuth2Github(c *gin.Context) {
	oauth2cfg := oauth2.Config{
		ClientID:     conf.Github.ClientID,
		ClientSecret: conf.Github.ClientSecret,
		Scopes:       []string{},
		Endpoint:     github.Endpoint,
	}
	url := oauth2cfg.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// OAuth2GithubCallback Github OAuth2 Callback
// Find or create connection, update user
func OAuth2GithubCallback(c *gin.Context) {
	oauth2cfg := oauth2.Config{
		ClientID:     conf.Github.ClientID,
		ClientSecret: conf.Github.ClientSecret,
		Scopes:       []string{},
		Endpoint:     github.Endpoint,
	}
	code := c.Query("code")

	ctx := context.Background()
	tok, err := oauth2cfg.Exchange(ctx, code)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	client := oauth2cfg.Client(ctx, tok)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer resp.Body.Close()

	g := githubProfile{}
	json.NewDecoder(resp.Body).Decode(&g)

	user := models.User{GithubID: g.ID}
	if db.DB.Where("github_id = ?", g.ID).First(&user).RecordNotFound() {
		db.DB.Transaction(func(t *gorm.DB) error {
			if err := db.DB.Create(&user).Error; err != nil {
				return err
			}

			id, err := uuid.NewDCESecurity(uuid.Person, uint32(user.ID))
			if err != nil {
				return err
			}

			user.Name = g.Name
			user.UUID = id.String()
			if err := db.DB.Save(&user).Error; err != nil {
				return err
			}

			return nil
		})
	}

	tokenString, err := token.AccessToken(user.UUID)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/login?token="+tokenString)
}
