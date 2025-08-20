package post

import (
	"time"

	"github.com/ardianilyas/go-feature-based/internal/auth"
	"github.com/ardianilyas/go-feature-based/internal/category"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID 			uuid.UUID 			`gorm:"type:uuid;primaryKey" json:"id"`
	Title 		string 				`gorm:"type:varchar(255); not null" json:"title"`
	Content 	string 				`gorm:"type:text;not null" json:"content"`
	UserID 		uuid.UUID 			`gorm:"type:uuid;not null" json:"-"`
	User 		auth.User 			`gorm:"foreignKey:UserID" json:"author"`
	CategoryID 	uuid.UUID 			`gorm:"type:uuid;not null" json:"-"`
	Category 	category.Category 	`gorm:"foreignKey:CategoryID" json:"category"`
	CreatedAt 	time.Time 			`json:"created_at"`
	UpdatedAt 	time.Time 			`json:"updated_at"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}