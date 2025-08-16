package auth

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100)" json:"name"`
	Email     string    `gorm:"uniqueIndex;type:varchar(100)" json:"email"`
	Password  string    `gorm:"type:varchar(255)" json:"-"`
	Role      string    `gorm:"default:user" json:"role"`
	CreatedAt time.Time
	UpdatedAt time.Time
}