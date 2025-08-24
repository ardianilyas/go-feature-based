package seed

import (
	"fmt"
	"log"
	"time"

	"github.com/ardianilyas/go-feature-based/config"
	"github.com/ardianilyas/go-feature-based/internal/auth"
	"github.com/ardianilyas/go-feature-based/internal/category"
	"github.com/ardianilyas/go-feature-based/internal/post"
	"github.com/brianvoe/gofakeit/v6"
	"gorm.io/gorm"
)

func SeedPost(n int) {
	config.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&post.Post{})

	var users []auth.User
	var categories []category.Category

	if err := config.DB.Find(&users).Error; err != nil {
		log.Fatalf("Error finding users: %v", err)
	}
	if err := config.DB.Find(&categories).Error; err != nil {
		log.Fatalf("Error finding categories: %v", err)
	}

	if len(users) == 0 || len(categories) == 0 {
		log.Println("No users or categories found. Please seed users and categories first.")
		return
	}

	for i := 0; i < n; i++ {
		post := post.Post{
			Title: gofakeit.Sentence(6),
			Content: gofakeit.Paragraph(3, 5, 12, ""),
			UserID: users[gofakeit.Number(0, len(users)-1)].ID,
			CategoryID: categories[gofakeit.Number(0, len(categories)-1)].ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := config.DB.Create(&post).Error; err != nil {
			fmt.Println("Error seeding post:", err)
		}
	}

	fmt.Printf("%d fake posts seeded\n", n)
}