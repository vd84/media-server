package data

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string `gorm:"type:varchar(100);uniqueIndex"`
	HashedPassword string
}
