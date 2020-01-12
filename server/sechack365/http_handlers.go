package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

type authHandler struct {
}

func (a *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.auth(w, r)
	default:
		respond(w, http.StatusNotFound)
	}
}

func (a *authHandler) auth(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie(cookieNameSessionID)
	response := map[string]bool{
		"authenticated": err != nil,
	}
	respondJSON(w, http.StatusOK, response)
}

type loginHandler struct {
}

func (l *loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		l.login(w, r)
	default:
		respond(w, http.StatusNotFound)
	}
}

func (l *loginHandler) login(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if !notEmptyAll(username, password) {
		respondMessage(w, http.StatusBadRequest, "EMPTY_PARAMETER")
		return
	}
	count := 0
	db.Model(&user{}).Where("name = ?", username).Count(&count)
	if count == 0 {
		respondMessage(w, http.StatusOK, "USERNAME_DOES_NOT_EXIST")
		return
	}
	u := user{}
	db.Where("name = ?", username).First(&u)
	if db.Error != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		respondMessage(w, http.StatusOK, "PASSWORD_DOES_NOT_MATCH")
		return
	}
	login(w, u)
	respondJSON(w, http.StatusOK, u)
}

type lessonsHandler struct {
}

func (l *lessonsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		l.fetch(w, r)
	case http.MethodPost:
		l.store(w, r)
	default:
		respond(w, http.StatusNotFound)
	}
}

func (l *lessonsHandler) fetch(w http.ResponseWriter, r *http.Request) {
	lessons := []lesson{}
	err := fetch(w, r, &lessons)
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, lessons)
}

func (l *lessonsHandler) store(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1 << 62)
	title := r.MultipartForm.Value["title"][0]
	description := r.MultipartForm.Value["description"][0]
	image := r.MultipartForm.Value["os"][0]
	consolePort := r.MultipartForm.Value["consolePort"][0]
	userID := r.MultipartForm.Value["userID"][0]
	if !notEmptyAll(title, description, image, consolePort, userID) {
		respond(w, http.StatusBadRequest)
		return
	}
	consolePortUint, err := strconv.ParseUint(consolePort, 10, 64)
	if err != nil {
		respond(w, http.StatusBadRequest)
		return
	}
	userIDUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		respond(w, http.StatusBadRequest)
		return
	}
	lesson := lesson{
		Title:       title,
		Description: description,
		Book:        "",
		ConsolePort: uint(consolePortUint),
		UserID:      uint(userIDUint),
	}
	tx := db.Begin()
	tx.Save(&lesson)
	if tx.Error != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	ports := r.MultipartForm.Value["ports[]"]
	for _, port := range ports {
		portUint, err := strconv.ParseUint(port, 10, 64)
		if err != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
		lessonPort := lessonPort{
			Port:     uint(portUint),
			LessonID: lesson.ID,
		}
		tx.Save(&lessonPort)
		if tx.Error != nil {
			tx.Rollback()
			respond(w, http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()
	u := user{}
	db.Where("id = ?", userID).First(&u)
	if db.Error != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	lessonDirectoryPath := filepath.Join(pathLessons, fmt.Sprintf("%d", u.ID), fmt.Sprintf("%d", lesson.ID))
	os.MkdirAll(lessonDirectoryPath, 0755)
	thumbnailHeaders := r.MultipartForm.File["thumbnail"]
	if len(thumbnailHeaders) == 1 {
		thumbnailHeader := thumbnailHeaders[0]
		filenameComponents := strings.Split(thumbnailHeader.Filename, ".")
		if 2 <= len(filenameComponents) {
			extension := filenameComponents[len(filenameComponents)-1]
			thumbnailFilePath := filepath.Join(lessonDirectoryPath, "thumbnail."+extension)
			thumbnailFile, err := os.Create(filepath.Join(thumbnailFilePath))
			if err != nil {
				respond(w, http.StatusInternalServerError)
				return
			}
			defer thumbnailFile.Close()
			thumbnail, err := thumbnailHeader.Open()
			if err != nil {
				respond(w, http.StatusInternalServerError)
				return
			}
			io.Copy(thumbnailFile, thumbnail)
		}
	}
	df := newDockerfile(image, u.Name)
	dockerfilePath := filepath.Join(lessonDirectoryPath, "Dockerfile")
	err = df.write(dockerfilePath)
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	imagename, err := uuid.NewUUID()
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	d := docker{}
	err = d.buildImage(imagename.String(), filepath.Dir(dockerfilePath))
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	containername, err := uuid.NewUUID()
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	containerID, err := d.runContainer(containername.String(), imagename.String(), ports...)
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	lesson.DockerContainerID = string(containerID)
	db.Save(&lesson)
	respondJSON(w, http.StatusOK, lesson)
}

type logoutHandler struct {
}

func (l *logoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		l.logout(w, r)
	default:
		respond(w, http.StatusNotFound)
	}
}

func (l *logoutHandler) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, newCookie(cookieNameSessionID, "", -1))
}

type registerHandler struct {
}

func (h *registerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.register(w, r)
	default:
		respond(w, http.StatusNotFound)
	}
}

func (h *registerHandler) register(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if !notEmptyAll(username, password) {
		respondMessage(w, http.StatusBadRequest, "EMPTY_PARAMETER")
		return
	}
	count := 0
	db.Model(&user{}).Where("name = ?", username).Count(&count)
	if count != 0 {
		respondMessage(w, http.StatusOK, "USERNAME_ALREADY_EXISTS")
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		respond(w, http.StatusInternalServerError)
	}
	u := user{Name: username, Password: string(hashedPassword)}
	db.Save(&u)
	if db.Error != nil {
		respond(w, http.StatusInternalServerError)
	}
	login(w, u)
	respondJSON(w, http.StatusOK, u)
}

type userHandler struct {
}

func (u *userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		u.fetch(w, r)
	default:
		respond(w, http.StatusNotFound)
	}
}

func (u *userHandler) fetch(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieNameSessionID)
	if err != nil {
		return
	}
	session, err := sessionManager.Provider.Get(cookie.Value)
	if err != nil {
		return
	}
	user := user{}
	err = session.Object(sessionKeyUser, &user)
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, user)
}

type usersHandler struct {
}

func (u *usersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		u.fetch(w, r)
	default:
		respond(w, http.StatusNotFound)
	}
}

func (u *usersHandler) fetch(w http.ResponseWriter, r *http.Request) {
	users := []user{}
	err := fetch(w, r, &users)
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, users)
}
