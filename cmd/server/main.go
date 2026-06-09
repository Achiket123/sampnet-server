package main

import (
	"server/internal/app/routes"
	"server/internal/platform/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load(".env")
	database.Init()

}

func main() {
	// Perform database migrations
	if err := database.PerformMigrations(); err != nil {
		panic("Failed to perform database migrations: " + err.Error())
	}
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
	}))
	routes.SetupRoutes(r)

	r.Run(":8000")
}
