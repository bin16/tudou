package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bin16/tudou/db"
	"github.com/bin16/tudou/models"
	"github.com/bin16/tudou/vars"
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
	}).Order("time_start, duration, id").Not("status", vars.TodoStatusRemoved).Find(&todayTodos)
	tomorrowTodos := []models.Todo{}
	db.DB.Debug().Preload("Event").Where(map[string]interface{}{
		"user_id": u.ID,
		"date":    tomorrow,
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
