package Mysql

import (
	"gorm.io/gorm"
	"time"
)

type Base struct {
	ID        string         `gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
