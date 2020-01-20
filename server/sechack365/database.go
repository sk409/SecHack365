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
	db.AutoMigrate(material{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").AddForeignKey("author_user_id", "users(id)", "NO ACTION", "NO ACTION")
	db.AutoMigrate(lessonMaterial{}).AddForeignKey("lesson_id", "lessons(id)", "CASCADE", "CASCADE").AddForeignKey("material_id", "materials(id)", "CASCADE", "CASCADE")
	db.AutoMigrate(follow{}).AddForeignKey("following_user_id", "users(id)", "CASCADE", "CASCADE").AddForeignKey("followed_user_id", "users(id)", "CASCADE", "CASCADE").AddUniqueIndex("following_user_id_followed_user_id_unique", "following_user_id", "followed_user_id")
	db.AutoMigrate(download{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").AddForeignKey("original_material_id", "materials(id)", "NO ACTION", "CASCADE").AddForeignKey("copied_material_id", "materials(id)", "CASCADE", "CASCADE")
}
