package controllers

import (
	"net/http"

	"github.com/bin16/tudou/db"
	"github.com/bin16/tudou/models"
	"github.com/bin16/tudou/vars"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// FinishTodo POST /a/todo.finish?id=X
func FinishTodo(c *gin.Context) {
	todoID := getTodoID(c)
	todo := models.Todo{ID: todoID}
	if err := db.DB.Preload("Event").Where("id = ?", todoID).First(&todo).Error; err != nil {
		c.JSON(msgJSON(http.StatusInternalServerError, err.Error()))
		return
	}
	user := c.MustGet("user").(*models.User)

	// only [today]'s [open] todo can be finished
	if todo.Status != vars.TodoStatusOpen || todo.Date != getToday() {
		c.Status(http.StatusBadRequest)
		return
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		// 1. Finish todo
		todo.Event.Status = vars.EventStatusClosed
		todo.Status = vars.TodoStatusFinished
		if err := db.DB.Save(&todo).Error; err != nil {
			return err
		}

		// 2. Change user's coins
		earnedCoins := calcCoins(vars.ActionFinishTodo)
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
			Action:  vars.ActionFinishTodo,
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

// CancelTodo POST /a/todo.cancel?id=X
func CancelTodo(c *gin.Context) {
	todoID := getTodoID(c)
	todo := models.Todo{ID: todoID}
	if err := db.DB.Preload("Event").Where("id = ?", todoID).First(&todo).Error; err != nil {
		c.JSON(msgJSON(http.StatusInternalServerError, err.Error()))
		return
	}
	user := c.MustGet("user").(*models.User)

	// only [open] todo can be canceled
	if todo.Status != vars.TodoStatusOpen {
		c.Status(http.StatusBadRequest)
		return
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		// 1. Cancel todo
		todo.Event.Status = vars.EventStatusClosed
		todo.Status = vars.TodoStatusCanceled
		err = tx.Save(&todo).Error
		if err != nil {
			return err
		}
		// 2. Change user's coins
		earnedCoins := calcCoins(vars.ActionCancelTodo)
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
			Action:  vars.ActionCancelTodo,
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

// PushTodo POST /a/todo.push?id=X
// Postpone the todo, today's todo is postponed,
// also generate open new one at tomorrow
func PushTodo(c *gin.Context) {
	todoID := getTodoID(c)
	todo := models.Todo{ID: todoID}
	if err := db.DB.Preload("Event").Where("id = ?", todoID).First(&todo).Error; err != nil {
		c.JSON(msgJSON(http.StatusInternalServerError, err.Error()))
		return
	}
	user := c.MustGet("user").(*models.User)

	// only today's open todos can be postponed
	if todo.Status != vars.TodoStatusOpen || todo.Date != getToday() {
		c.Status(http.StatusBadRequest)
		return
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		// 1. Postpone todo
		todo.Status = vars.TodoStatusPostponed
		err = tx.Save(&todo).Error
		if err != nil {
			return err
		}
		// 2. Copy todo
		copiedTodo := models.Todo{
			EventID: todo.EventID,
			UserID:  user.ID,
			Date:    getTomorraw(),
			Status:  vars.TodoStatusOpen,
		}
		err = tx.Create(&copiedTodo).Error
		if err != nil {
			return err
		}
		// 3. Change user's coins
		earnedCoins := calcCoins(vars.ActionPushTodo)
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
			Action:  vars.ActionPushTodo,
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

// PullTodo POST /a/todo.pull?id=X
// Pull tomorrow's todo item
func PullTodo(c *gin.Context) {
	todoID := getTodoID(c)
	todo := models.Todo{ID: todoID}
	if err := db.DB.Preload("Event").First(&todo).Error; err != nil {
		c.JSON(msgJSON(http.StatusInternalServerError, err.Error()))
		return
	}
	user := c.MustGet("user").(*models.User)

	today := getToday()
	// pull tomorraw's open todos only
	if todo.Date == today || todo.Status != vars.TodoStatusOpen {
		c.Status(http.StatusBadRequest)
		return
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		// 1. Pull todo
		todo.Date = today
		todo.Status = vars.TodoStatusOpen
		err = tx.Save(&todo).Error
		if err != nil {
			return err
		}
		// 2. Change user's coins
		earnedCoins := calcCoins(vars.ActionPullTodo)
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
			Action:  vars.ActionPullTodo,
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

// RemoveTodo POST /a/todo.remove?id=X
// remove event and all related todos
func RemoveTodo(c *gin.Context) {
	todoID := getTodoID(c)
	todo := models.Todo{ID: todoID}
	if err := db.DB.Preload("Event").First(&todo).Error; err != nil {
		c.JSON(msgJSON(http.StatusInternalServerError, err.Error()))
		return
	}
	user := c.MustGet("user").(*models.User)

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		// 1. Remove todos
		err = tx.Where("event_id = ?", todo.EventID).Delete(&models.Todo{}).Error
		if err != nil {
			return err
		}
		// 2. Remove event
		err = tx.Delete(&todo.Event).Error
		if err != nil {
			return err
		}
		// 3. Change user's coins
		earnedCoins := calcCoins(vars.ActionRemoveTodo)
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
			Action:  vars.ActionRemoveTodo,
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
