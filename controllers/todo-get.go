package controllers

import (
	"net/http"

	"github.com/bin16/tudou/vars"

	"github.com/bin16/tudou/db"
	"github.com/bin16/tudou/models"
	"github.com/gin-gonic/gin"
)

// GetTodos GET /api/my/todos
func GetTodos(c *gin.Context) {
	u := c.MustGet("user").(*models.User)
	today := getToday()
	tomorrow := getTomorrow()
	todayTodos := []models.Todo{}
	db.DB.Preload("Event").Where(map[string]interface{}{
		"user_id": u.ID,
		"date":    today,
		// "status":  vars.TodoStatusOpen,
	}).Order("time_start, duration, id").Find(&todayTodos)
	tomorrowTodos := []models.Todo{}
	db.DB.Debug().Preload("Event").Where(map[string]interface{}{
		"user_id": u.ID,
		"date":    tomorrow,
		// "status":  vars.TodoStatusOpen,
	}).Order("time_start, duration, id").Find(&tomorrowTodos)
	overdueTodos := []models.Todo{}
	db.DB.Preload("Event").Where(map[string]interface{}{
		"user_id": u.ID,
		"status":  vars.TodoStatusOpen,
	}).Order("time_start, duration, id").Not("date", []string{today, tomorrow}).Find(&overdueTodos)

	c.JSON(http.StatusOK, gin.H{
		"today":    today,
		"tomorrow": tomorrow,
		today:      todayTodos,
		tomorrow:   tomorrowTodos,
		"overdue":  overdueTodos,
	})
}
