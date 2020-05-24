package main

import (
	"fmt"
	"log"

	"github.com/bin16/tudou/conf"
	"github.com/bin16/tudou/controllers"
	"github.com/bin16/tudou/db"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

// tudou config.toml
func main() {
	conf.Load()

	database, dbERR := gorm.Open(conf.Database.Dialect, conf.Database.DB)
	if dbERR != nil {
		log.Fatal(dbERR)
	}
	db.Use(database)

	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("static", false)))
	r.GET("/", controllers.IndexPage)
	r.GET("/note", controllers.NotePage)
	r.GET("/login", controllers.LoginPage)
	if conf.Github.Enabled {
		r.GET("/auth/github", controllers.OAuth2Github)
		r.GET("/auth/github/callback", controllers.OAuth2GithubCallback)
	}

	api := r.Group("/api")
	api.Use(controllers.ShouldAuthed)
	{
		api.GET("/my/profile", controllers.GetProfile)
		api.GET("/my/todos", controllers.GetTodos)
		api.GET("/my/logs", controllers.GetLogs)
		api.GET("/my/notes", controllers.GetNotes)
		api.GET("/my/settings", controllers.GetSettings)
	}

	actions := api.Group("/a")
	actions.Use(controllers.ShouldAuthed)
	{
		actions.POST("/todo.update", controllers.UpdateTodo)
		// Tomorraw
		actions.POST("/todo.create", controllers.CreateTodo)
		actions.POST("/todo.clone", controllers.CloneTodo)
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
		// File
		actions.POST("/file.upload", controllers.UploadFile)
	}

	r.Run(fmt.Sprintf(":%d", conf.Server.Port))
}
