package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bin16/tudou/vars"
	"github.com/gin-gonic/gin"
)

func http2Push(c *gin.Context, urls ...string) {
	if pusher := c.Writer.Pusher(); pusher != nil {
		for _, u := range urls {
			pusher.Push(u, nil) // ignore error
		}
	}
}

func calcCoins(action string) int64 {
	var coin int64 = 5
	switch action {
	case vars.ActionFinishTodo:
		return 1 * coin
	case vars.ActionPushTodo:
		return -1 * coin
	case vars.ActionCancelTodo:
		return -1 * coin
	case vars.ActionRemoveTodo:
		return -10 * coin
	case vars.ActionCreateTodo:
		return 1 * coin
	case vars.ActionRenewTodo:
		return -20 * coin
	}

	return 0
}

func errorJSON(err error) (int, gin.H) {
	if err == nil {
		return http.StatusOK, gin.H{
			"ok": 1,
		}
	}

	return http.StatusInternalServerError, gin.H{
		"ok":      0,
		"message": err.Error(),
	}
}

func getTodoID(c *gin.Context) int {
	todoID := c.Query("id")
	id, _ := strconv.Atoi(todoID)

	return id
}

func getToday() string {
	return time.Now().Format("2006-01-02")
}

func getTomorraw() string {
	return time.Now().AddDate(0, 0, 1).Format("2006-01-02")
}
