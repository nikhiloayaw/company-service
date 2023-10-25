package domain

import (
	"auth-service/pkg/utils"

	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey;not null"`
	ID26     string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

func (b *User) BeforeCreate(tx *gorm.DB) (err error) {

	b.ID26 = utils.GenerateUUID()

	return nil
}
