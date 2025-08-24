package seed

import (
	"github.com/ardianilyas/go-feature-based/config"
)

func ResetTables() {
	config.DB.Exec("TRUNCATE TABLE categories, posts, users RESTART IDENTITY CASCADE")
}

func RunSeeders() {
	ResetTables()
	SeedUser()
	SeedCategory(50)
	SeedPost(40)
}