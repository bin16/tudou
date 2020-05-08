package main

import (
	"log"

	"github.com/bin16/tudou/controllers"
	"github.com/bin16/tudou/db"
	"github.com/bin16/tudou/models"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	if dbERR := db.Open(); dbERR != nil {
		log.Fatal(dbERR)
	}
	defer db.DB.Close()
	models.AutoMigrate(db.DB)

	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("static", false)))
	r.GET("/", controllers.IndexPage)
	r.GET("/lab", controllers.LabPage)
	r.GET("/note", controllers.NotePage)

	api := r.Group("/api", controllers.ShouldAuthed)
	{
		api.GET("/my/profile", controllers.GetProfile)
		api.GET("/my/todos", controllers.GetTodos)
		api.GET("/my/logs", controllers.GetLogs)
		api.GET("/my/notes", controllers.GetNotes)
	}

	actions := api.Group("/a")
	{
		actions.POST("/user.create", controllers.CreateUser)
		actions.POST("/todo.update", controllers.UpdateTodo)
		// Tomorraw
		actions.POST("/todo.create", controllers.CreateTodo)
		// Today
		actions.POST("/todo.finish", controllers.FinishTodo)
		actions.POST("/todo.cancel", controllers.CancelTodo)
		actions.POST("/todo.push", controllers.PushTodo)
		actions.POST("/todo.pull", controllers.PullTodo)
		// Overdue
		actions.POST("/todo.renew", controllers.RenewTodo)
		// Other
		actions.POST("/todo.remove", controllers.RemoveTodo)
		// ...
		actions.POST("/settings.update", controllers.UpdateSettings)
		// Notes
		actions.POST("/note.create", controllers.CreateNote)
		actions.POST("/note.update", controllers.UpdateNote)
		actions.POST("/note.remove", controllers.RemoveNote)
	}

	r.Run(":2358")
}
