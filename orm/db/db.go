package db

import (
	"gorm.io/gorm"
)

type GormDB interface {
	GetGormDB() *gorm.DB
}
