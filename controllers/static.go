package controllers

import (
	"github.com/gin-gonic/gin"
)

// IndexPage index.html
func IndexPage(c *gin.Context) {
	c.File("static/index.html")
}

// LabPage lab.html
func LabPage(c *gin.Context) {
	c.File("static/lab.html")
}

// NotePage note.html
func NotePage(c *gin.Context) {
	c.File("static/note.html")
}
