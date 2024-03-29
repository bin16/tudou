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

// RenewTodo POST /a/todo.renew
func RenewTodo(c *gin.Context) {
	todoID := getTodoID(c)
	todo := models.Todo{ID: todoID}
	if err := db.DB.Preload("Event").Where("id = ?", todoID).First(&todo).Error; err != nil {
		c.JSON(msgJSON(http.StatusInternalServerError, err.Error()))
		return
	}
	user := c.MustGet("user").(*models.User)
	tz := user.Setting.Timezone

	if todo.Status != vars.TodoStatusOpen || todo.Date == getToday(tz) || todo.Date == getTomorrow(tz) {
		c.Status(http.StatusBadRequest)
		return
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		// 1. Renew todo
		todo.Status = vars.TodoStatusOverdue
		err = tx.Save(&todo).Error
		if err != nil {
			return err
		}
		// 2. Copy todo
		copiedTodo := models.Todo{
			EventID:   todo.EventID,
			UserID:    user.ID,
			Date:      getTomorrow(user.Setting.Timezone),
			Status:    vars.TodoStatusOpen,
			TimeStart: todo.TimeStart,
			Duration:  todo.Duration,
			TimeSet:   todo.TimeSet,
		}
		err = tx.Create(&copiedTodo).Error
		if err != nil {
			return err
		}
		// 3. Change user's coins
		earnedCoins := calcCoins(vars.ActionRenewTodo)
		user.Coins += earnedCoins
		err = tx.Save(user).Error
		if err != nil {
			return err
		}
		// 4. log
		err = tx.Create(&models.Log{
			TodoID:  todo.ID,
			EventID: todo.EventID,
			UserID:  user.ID,
			Action:  vars.ActionRenewTodo,
			Coins:   earnedCoins,
		}).Error
		if err != nil {
			return err
		}
		err = tx.Create(&models.Log{
			TodoID:  copiedTodo.ID,
			EventID: todo.EventID,
			UserID:  user.ID,
			Action:  vars.ActionCopyTodo,
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

	http2Push(c, "/api/my/todos")

	c.Status(http.StatusOK)
}

// CreateTodo POST /a/todo.create
func CreateTodo(c *gin.Context) {
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

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		// 1. Create todo
		todo := models.Todo{
			Event: models.Event{
				Content: values.Content,
				Status:  vars.EventStatusOpen,
				UserID:  user.ID,
			},
			User:      *user,
			Status:    vars.TodoStatusOpen,
			UserID:    user.ID,
			Date:      getTomorrow(user.Setting.Timezone),
			TimeStart: values.TimeStart(),
			Duration:  values.Duration,
			TimeSet:   values.TimeSet,
		}
		if err := db.DB.Create(&todo).Related(&todo.Event).Error; err != nil {
			return err
		}
		// 2. Change user's coins
		earnedCoins := calcCoins(vars.ActionCreateTodo)
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
			Action:  vars.ActionCreateTodo,
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

// CloneTodo POST /a/todo.clone
func CloneTodo(c *gin.Context) {
	todoID := getTodoID(c)

	todo := models.Todo{ID: todoID}
	if err := db.DB.Preload("Event").First(&todo).Error; err != nil {
		c.JSON(msgJSON(http.StatusInternalServerError, err.Error()))
		return
	}

	user := c.MustGet("user").(*models.User)

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		// 1. Create todo & event
		todo := models.Todo{
			Event: models.Event{
				Content: todo.Event.Content,
				Status:  vars.EventStatusOpen,
				UserID:  user.ID,
			},
			User:      *user,
			Status:    vars.TodoStatusOpen,
			UserID:    user.ID,
			Date:      getTomorrow(user.Setting.Timezone),
			TimeStart: todo.TimeStart,
			Duration:  todo.Duration,
			TimeSet:   todo.TimeSet,
		}
		if err := db.DB.Create(&todo).Related(&todo.Event).Error; err != nil {
			return err
		}
		// 2. Change user's coins
		earnedCoins := calcCoins(vars.ActionCloneTodo)
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
			Action:  vars.ActionCloneTodo,
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

	c.Status(http.StatusOK)
}
