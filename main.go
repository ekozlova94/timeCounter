package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"log"
	"timeCounter/handlers"
)

func main() {
	var config = flag.String("config-dir", ".", "path to config file")
	flag.Parse()
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(*config)  // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	r := gin.Default()
	r.POST("/api/start", handlers.Start)
	r.POST("/api/stop", handlers.Stop)
	r.POST("/api/edit", handlers.Edit)
	r.GET("/api/info", handlers.Info)
	r.GET("/api/today", handlers.Today)
	r.POST("/api/start-break", handlers.BreakStart)
	r.POST("/api/stop-break", handlers.BreakStop)
	r.POST("/api/edit-break", handlers.EditBreak)

	port := viper.GetString("server.port")
	if port == "" {
		port = "6000"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
