package controllers

import (
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	file, _ := c.FormFile("file")
	fullpath := path.Join("/uploads", file.Filename)
	c.SaveUploadedFile(file, path.Join("static", fullpath))
	c.JSON(http.StatusCreated, gin.H{
		"url": fullpath,
	})
}
