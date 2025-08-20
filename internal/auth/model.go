package auth

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100)" json:"name"`
	Email     string    `gorm:"uniqueIndex;type:varchar(100)" json:"email"`
	Password  string    `gorm:"type:varchar(255)" json:"-"`
	Role      string    `gorm:"default:user" json:"role"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}