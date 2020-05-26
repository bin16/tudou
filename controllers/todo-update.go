package controllers

import (
	"net/http"

	"github.com/bin16/tudou/controllers/forms"
	"github.com/bin16/tudou/db"
	"github.com/bin16/tudou/models"
	"github.com/bin16/tudou/vars"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// UpdateTodo POST /a/todo.update?id=X
func UpdateTodo(c *gin.Context) {
	values := forms.TodoForm{}
	if err := c.BindJSON(&values); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if valid, message := values.Valid(); !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	user := c.MustGet("user").(*models.User)

	todo := models.Todo{}
	if err := db.DB.Preload("Event").Where(map[string]interface{}{
		"id":      values.ID,
		"user_id": user.ID,
	}).First(&todo).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	if todo.Status != vars.TodoStatusOpen || todo.Event.Status != vars.EventStatusOpen {
		c.Status(http.StatusBadRequest)
		return
	}
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		// 1. Update todo
		todo.Event.Content = values.Content
		todo.TimeStart = values.TimeStart()
		todo.Duration = values.Duration
		err = db.DB.Save(&todo).Error
		if err != nil {
			return err
		}

		// 2. Change user's coins
		earnedCoins := calcCoins(vars.ActionEditTodo)
		user.Coins += earnedCoins
		err = tx.Save(user).Error
		if err != nil {
			return err
		}
		// 3. log
		err = tx.Create(&models.Log{
			TodoID:  todo.ID,
			EventID: todo.EventID,
			UserID:  user.ID,
			Action:  vars.ActionEditTodo,
			Coins:   earnedCoins,
		}).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		c.JSON(msgJSON(http.StatusInternalServerError, err.Error()))
		return
	}

	c.Status(http.StatusCreated)
}
