package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bin16/tudou/vars"
	"github.com/gin-gonic/gin"
)

func http2Push(c *gin.Context, urls ...string) {
	log.Println("http2Push ====")
	authorization := c.GetHeader("Authorization")
	if pusher := c.Writer.Pusher(); pusher != nil {
		for _, u := range urls {
			err := pusher.Push(u, &http.PushOptions{
				Header: http.Header{
					"Authorization": []string{authorization},
					"From":          []string{c.Request.URL.Path},
				},
			})
			if err != nil {
				log.Printf("---- Failed to push %s\n", u)
			}
		}
	} else {
		log.Println("==== Pusher not available")
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
		return 0 * coin
	case vars.ActionCreateTodo:
		return 1 * coin
	case vars.ActionCloneTodo:
		return 1 * coin
	case vars.ActionRenewTodo:
		return -20 * coin
	}

	return 0
}

func msgJSON(status int, message string) (int, gin.H) {
	return status, gin.H{"message": message}
}

func getTodoID(c *gin.Context) int {
	todoID := c.Query("id")
	id, _ := strconv.Atoi(todoID)

	return id
}

func getToday(timezone string) string {
	t := time.Now()
	local, err := time.LoadLocation(timezone)
	if err != nil {
		return t.UTC().Format("2006-01-02")
	}
	localtime := t.In(local)

	return localtime.Format("2006-01-02")
}

func getTomorrow(timezone string) string {
	t := time.Now().AddDate(0, 0, 1)
	local, err := time.LoadLocation(timezone)
	if err != nil {
		return t.UTC().Format("2006-01-02")
	}
	localtime := t.In(local)

	return localtime.Format("2006-01-02")
}
