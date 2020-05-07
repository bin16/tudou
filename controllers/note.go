package controllers

import (
	"net/http"

	"github.com/bin16/tudou/controllers/forms"
	"github.com/bin16/tudou/db"
	"github.com/bin16/tudou/models"
	"github.com/gin-gonic/gin"
)

func GetNotes(c *gin.Context) {
	notes := []models.Note{}
	db.DB.Where(map[string]interface{}{
		"archived": false,
	}).Find(&notes)

	c.JSON(http.StatusOK, gin.H{
		"notes": notes,
	})
}

func CreateNote(c *gin.Context) {
	values := forms.NoteForm{}
	if err := c.BindJSON(&values); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if valid, _ := values.Valid(); !valid {
		c.Status(http.StatusBadRequest)
		return
	}
	user := c.MustGet("user").(*models.User)

	note := models.Note{
		Content:  values.Content,
		Title:    values.Title,
		UserID:   user.ID,
		User:     *user,
		Archived: false,
	}
	if err := db.DB.Create(&note).Error; err != nil {
		c.JSON(errorJSON(err))
		return
	}

	c.Status(http.StatusCreated)
}

func UpdateNote(c *gin.Context) {
	values := forms.NoteForm{}
	if err := c.BindJSON(&values); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if valid, message := values.Valid(); !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}
	user := c.MustGet("user").(*models.User)

	note := models.Note{}
	if err := db.DB.Debug().Where(map[string]interface{}{
		"id":      values.ID,
		"user_id": user.ID,
	}).First(&note).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	note.Content = values.Content
	note.Title = values.Title
	if err := db.DB.Save(&note).Error; err != nil {
		c.JSON(errorJSON(err))
		return
	}

	c.Status(http.StatusOK)
}