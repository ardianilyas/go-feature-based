package migrations

import (
	"log"

	"github.com/ardianilyas/go-feature-based/config"
	"github.com/ardianilyas/go-feature-based/internal/auth"
	"github.com/ardianilyas/go-feature-based/internal/category"
)

func RunMigrations() {
	err := config.DB.AutoMigrate(&auth.User{}, &category.Category{})
	if err != nil {
		log.Fatal("Error running migrations:", err)
	}

	log.Println("Migrations completed")
}