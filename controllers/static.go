package controllers

import (
	"path"

	"github.com/bin16/tudou/conf"
	"github.com/gin-gonic/gin"
)

// IndexPage index.html
func IndexPage(c *gin.Context) {
	c.File(path.Join(conf.Server.Static, "index.html"))
}

// NotePage note.html
func NotePage(c *gin.Context) {
	c.File(path.Join(conf.Server.Static, "note.html"))
}

// LoginPage note.html
func LoginPage(c *gin.Context) {
	c.File(path.Join(conf.Server.Static, "login.html"))
}
