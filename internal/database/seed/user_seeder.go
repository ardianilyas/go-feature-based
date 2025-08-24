package seed

import (
	"fmt"

	"github.com/ardianilyas/go-feature-based/config"
	"github.com/ardianilyas/go-feature-based/internal/auth"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUser() {
	config.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&auth.User{})

	password, _ := bcrypt.GenerateFromPassword([]byte("developer"), bcrypt.DefaultCost)

	user := auth.User{
		Name: "Ardian Ilyas",
		Email: "ardian@developer.com",
		Password: string(password),
		Role: "admin",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		fmt.Println("Error seeding user:", err)
		return
	}

	fmt.Println("User seeded successfully")
}