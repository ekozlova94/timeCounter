package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"timeCounter/handlers"
)

func main() {
	r := gin.Default()
	r.POST("/api/start", handlers.Start)
	r.POST("/api/stop", handlers.Stop)
	r.POST("/api/edit", handlers.Edit)
	r.GET("/api/info", handlers.Info)
	r.GET("/api/today", handlers.Today)
	r.POST("/api/start-break", handlers.BreakStart)
	r.POST("/api/stop-break", handlers.BreakStop)
	r.POST("/api/edit-break", handlers.EditBreak)
	if err := r.Run(); err != nil { // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
		log.Fatal(err)
	}
}
