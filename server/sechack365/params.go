package main

import (
	"os"
	"path/filepath"

	"github.com/sk409/gosession"
)

const (
	cookieMaxAge30Days  = 60 * 60 * 24 * 30
	cookieNameSessionID = "SESSION_ID"
	sessionKeyUser      = "SESSION_KEY_USER"
)

var (
	cwd                             string
	pathLessons                     string
	pathNoimage                     string
	pathPublic                      string
	pathPublicImages                string
	pathPublicImagesLessons         string
	pathPublicImagesLessonsClones   string
	pathPublicImagesMaterials       string
	pathPublicImagesMaterialsClones string
	pathPublicImagesUsers           string
	sessionManager                  = gosession.Manager{Provider: gosession.NewMemoryProvider()}
)

func init() {
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	pathLessons = filepath.Join(cwd, "lessons")
	pathPublic = filepath.Join(cwd, "public")
	pathPublicImages = filepath.Join(pathPublic, "images")
	pathPublicImagesLessons = filepath.Join(pathPublicImages, "lessons")
	pathPublicImagesMaterials = filepath.Join(pathPublicImages, "materials")
	pathPublicImagesUsers = filepath.Join(pathPublicImages, "users")
	pathNoimage = filepath.Join(pathPublic, "noimage.png")
	os.Mkdir(pathLessons, 0755)
	os.MkdirAll(pathPublicImagesLessons, 0755)
	os.MkdirAll(pathPublicImagesMaterials, 0755)
	os.MkdirAll(pathPublicImagesUsers, 0755)
}
