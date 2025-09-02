package main

import (
	"github.com/ardianilyas/go-feature-based/config"
	"github.com/ardianilyas/go-feature-based/internal"
	"github.com/ardianilyas/go-feature-based/internal/migrations"
	"github.com/ardianilyas/go-feature-based/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitLogger()
	config.LoadEnv()
	config.ConnectDB()
	migrations.RunMigrations()
	
	go middlewares.CleanupClients()

	r := gin.Default()

	r.Use(middlewares.CORSMiddleware())
	r.Use(middlewares.RateLimitMiddleware())
	
	internal.SetupRoutes(r)

	r.Run(":8000")
}