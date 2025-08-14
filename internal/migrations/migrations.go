package migrations

import (
	"log"

	"github.com/ardianilyas/go-feature-based/config"
	"github.com/ardianilyas/go-feature-based/internal/auth"
)

func RunMigrations() {
	err := config.DB.AutoMigrate(&auth.User{})
	if err != nil {
		log.Fatal("Error running migrations:", err)
	}

	log.Println("Migrations completed")
}