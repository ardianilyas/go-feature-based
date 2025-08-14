package main

import (
	"github.com/ardianilyas/go-feature-based/config"
	"github.com/ardianilyas/go-feature-based/internal/auth"
	"github.com/ardianilyas/go-feature-based/internal/migrations"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()
	migrations.RunMigrations()

	r := gin.Default()
	
	authRepo := auth.NewRepository()
	authService := auth.NewService(authRepo)
	auth.RegisterRoutes(r, authService)

	r.Run(":8000")
}