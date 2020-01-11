package main

import (
	"github.com/jinzhu/gorm"
)

type user struct {
	gorm.Model
	Name     string `gorm:"type:varchar(128);not null;unique"`
	Password string `gorm:"type:varchar(128);not null"`
}
