package main

import (
	"github.com/ardianilyas/go-feature-based/config"
	"github.com/ardianilyas/go-feature-based/internal/database/seed"
	"github.com/ardianilyas/go-feature-based/internal/migrations"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()
	migrations.RunMigrations()

	seed.RunSeeders()
}