package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username  string    `gorm:"unique;not null"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	RoleID    uuid.UUID `gorm:"type:uuid;not null"`
	Role      Role      `gorm:"foreignKey:RoleID"`
	Active    bool      `gorm:"default:true"`
	LastLogin time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) HashPassword() error {
	// Use bcrypt to hash the password
	// Param 1: password need to be hashed
	// Param 2: the cost of hashing, the higher the cost, the more secure the password, but the slower the hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	// Replace the password with the hashed version
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
