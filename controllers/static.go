package controllers

import (
	"github.com/gin-gonic/gin"
)

// IndexPage index.html
func IndexPage(c *gin.Context) {
	c.File("static/index.html")
}

// NotePage note.html
func NotePage(c *gin.Context) {
	c.File("static/note.html")
}

// LoginPage note.html
func LoginPage(c *gin.Context) {
	c.File("static/login.html")
}
