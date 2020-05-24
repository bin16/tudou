package controllers

import (
	"net/http"
	"strconv"

	"github.com/bin16/tudou/controllers/forms"
	"github.com/bin16/tudou/db"
	"github.com/bin16/tudou/models"
	"github.com/gin-gonic/gin"
)

func GetNotes(c *gin.Context) {
	notes := []models.Note{}
	db.DB.Preload("User").Where(map[string]interface{}{
		"archived": false,
	}).Order("created_at DESC").Find(&notes)

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
		Content:     values.Content,
		Title:       values.Title,
		UserID:      user.ID,
		User:        *user,
		Archived:    false,
		Attachments: values.Attachments,
	}
	if err := db.DB.Create(&note).Error; err != nil {
		c.JSON(msgJSON(http.StatusInternalServerError, err.Error()))
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
	if err := db.DB.Where(map[string]interface{}{
		"id":      values.ID,
		"user_id": user.ID,
	}).First(&note).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	note.Content = values.Content
	note.Title = values.Title
	note.Attachments = values.Attachments
	if err := db.DB.Save(&note).Error; err != nil {
		c.JSON(msgJSON(http.StatusInternalServerError, err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

func RemoveNote(c *gin.Context) {
	noteID := getNoteID(c)
	user := c.MustGet("user").(*models.User)

	err := db.DB.Delete(models.Note{ID: noteID, UserID: user.ID}).Error
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func getNoteID(c *gin.Context) int {
	id := c.Query("id")
	noteID, _ := strconv.Atoi(id)

	return noteID
}
