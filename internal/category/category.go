package category

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID 		uuid.UUID 	`gorm:"type:uuid;primaryKey" json:"id"`
	Name 	string 		`gorm:"type:varchar(100)" json:"name"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}