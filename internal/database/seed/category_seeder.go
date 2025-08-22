package seed

import (
	"fmt"

	"github.com/ardianilyas/go-feature-based/config"
	"github.com/ardianilyas/go-feature-based/internal/category"
	"github.com/brianvoe/gofakeit/v6"
	"gorm.io/gorm"
)

func SeedCategory(n int) {
	config.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&category.Category{})

	for i := 0; i < n; i++ {
		category := category.Category{
			Name: gofakeit.Word(),
		}

		if err := config.DB.Create(&category).Error; err != nil {
			fmt.Println("Error seeding category:", err)
		}
	}

	fmt.Printf("%d fake categories seeded", n)
}