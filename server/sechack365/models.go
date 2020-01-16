package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

type dockerContainer struct {
	id        string
	name      string
	hostPorts map[string]uint
}

type download struct {
	gorm.Model
	UserID     uint
	MaterialID uint
}

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

type follow struct {
	gorm.Model
	FollowingUserID uint `gorm:"not null"`
	FollowedUserID  uint `gorm:"not null"`
}

type user struct {
	gorm.Model
	Name             string `gorm:"type:varchar(128);not null;unique"`
	Password         string `gorm:"type:char(60);not null"`
	ProfileImagePath string `gorm:"type:varchar(256);"`
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
	ThumbnailPath     string     `gorm:"type:varchar(256);not null"`
	ConsolePort       uint       `gorm:"not null"`
	HostConsolePort   uint
	Downloaded        bool `gorm:"not null"`
	UserID            uint `gorm:"not null"`
}

func (l lesson) public() (interface{}, error) {
	p := structToMap(l)
	ports := []lessonPort{}
	db.Where("lesson_id = ?", l.ID).Find(&ports)
	if db.Error != nil {
		return nil, db.Error
	}
	p["Ports"] = ports
	return p, nil
}

type lessonMaterial struct {
	gorm.Model
	LessonID   uint `gorm:"not null"`
	MaterialID uint `gorm:"not null"`
}

func (lessonMaterial) TableName() string {
	return "lesson_material"
}

type lessonPort struct {
	gorm.Model
	Port     uint `gorm:"not null"`
	HostPort uint
	LessonID uint `gorm:"not null"`
}

type material struct {
	ID            uint `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `sql:"index"`
	Title         string     `gorm:"type:varchar(128);not null;"`
	Description   string     `gorm:"type:varchar(1024);not null;"`
	ThumbnailPath string     `gorm:"type:varchar(256);not null"`
	Downloaded    bool       `gorm:"not null"`
	UserID        uint       `gorm:"not null"`
	AuthorUserID  *uint
}

func (m material) public() (interface{}, error) {
	result := structToMap(m)
	lessonMaterials := []lessonMaterial{}
	db.Where("material_id = ?", m.ID).Find(&lessonMaterials)
	if db.Error != nil {
		return nil, db.Error
	}
	lessons := []lesson{}
	for _, lessonMaterial := range lessonMaterials {
		l := lesson{}
		db.Where("id = ?", lessonMaterial.LessonID).First(&l)
		lessons = append(lessons, l)
	}
	result["lessons"] = lessons
	return result, nil
}
