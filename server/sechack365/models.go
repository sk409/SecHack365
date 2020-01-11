package main

import (
	"github.com/jinzhu/gorm"
)

type user struct {
	gorm.Model
	Name     string `gorm:"type:varchar(128);not null;unique"`
	Password string `gorm:"type:char(60);not null"`
}

type lesson struct {
	gorm.Model
	Title             string `gorm:"type:varchar(256);not null"`
	Description       string `gorm:"type:text;not null"`
	Book              string `gorm:"type:text;not null"`
	DockerContainerID string `gorm:"type:char(64);"`
	ConsolePort       uint   `gorm:"not null"`
	UserID            uint   `gorm:"not null"`
}

type lessonPort struct {
	gorm.Model
	Port     uint `gorm:"not null"`
	LessonID uint `gorm:"not null"`
}
