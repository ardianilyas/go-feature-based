package main

import (
	"github.com/ardianilyas/go-feature-based/config"
	"github.com/ardianilyas/go-feature-based/internal"
	"github.com/ardianilyas/go-feature-based/internal/migrations"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()
	migrations.RunMigrations()

	r := gin.Default()
	
	internal.SetupRoutes(r)

	r.Run(":8000")
}