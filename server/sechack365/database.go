package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:root@(database)/sechack365?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(user{})
	db.AutoMigrate(lesson{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.AutoMigrate(lessonPort{}).AddForeignKey("lesson_id", "lessons(id)", "CASCADE", "CASCADE")
}
