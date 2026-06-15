package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"server/internal/app/routes"
	"server/internal/platform/cloudinary"
	"server/internal/platform/database"
	"server/internal/platform/redis"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load(".env")
	database.Init()
	redis.Init()
	cloudinary.Init()
}

func main() {
	// Listen for OS signals for graceful context cancellation
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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
	routes.SetupRoutes(ctx, r)

	r.Run(":8000")
}
