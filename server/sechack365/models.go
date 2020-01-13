package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

type file struct {
	Dir  string
	Name string
	Path string
	Text string
}

type folder struct {
	Dir          string
	Name         string
	Path         string
	ChildFiles   []file
	ChildFolders []folder
}

type user struct {
	gorm.Model
	Name     string `gorm:"type:varchar(128);not null;unique"`
	Password string `gorm:"type:char(60);not null"`
}

type lesson struct {
	ID                uint `gorm:"primary_key"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time `sql:"index"`
	Title             string     `gorm:"type:varchar(256);not null"`
	Description       string     `gorm:"type:text;not null"`
	Book              string     `gorm:"type:text;not null"`
	DockerContainerID string     `gorm:"type:char(64);"`
	ThumbnailPath     string     `gorm:"type:varchar(256)"`
	ConsolePort       uint       `gorm:"not null"`
	HostConsolePort   uint
	UserID            uint `gorm:"not null"`
}

func (l lesson) Public() (interface{}, error) {
	p := structToMap(l)
	ports := []lessonPort{}
	db.Where("lesson_id = ?", l.ID).Find(&ports)
	if db.Error != nil {
		return nil, db.Error
	}
	p["Ports"] = ports
	return p, nil
}

type lessonPort struct {
	gorm.Model
	Port     uint `gorm:"not null"`
	HostPort uint
	LessonID uint `gorm:"not null"`
}
