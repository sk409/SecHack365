package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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

type foldersHandler struct {
}

func (f *foldersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		f.fetch(w, r)
	default:
		respond(w, http.StatusNotFound)
	}
}

func (f *foldersHandler) fetch(w http.ResponseWriter, r *http.Request) {
	lessonID := r.URL.Query().Get("lessonID")
	path := r.URL.Query().Get("path")
	if !notEmptyAll(lessonID, path) {
		respond(w, http.StatusBadRequest)
		return
	}
	l := lesson{}
	db.Where("id = ?", lessonID).First(&l)
	d := docker{}
	output, err := d.exec(l.DockerContainerID, []string{}, "ls", "-al", path)
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	root := folder{Path: path}
	for _, line := range strings.Split(string(output), "\n") {
		if len(line) == 0 {
			continue
		}
		components := strings.Split(line, " ")
		name := components[len(components)-1]
		if name == "." || name == ".." {
			continue
		}
		switch line[0] {
		case '-':
			child := file{Dir: path, Name: name, Path: filepath.Join(path, name)}
			root.ChildFiles = append(root.ChildFiles, child)
		case 'd':
			child := folder{Dir: path, Name: name, Path: filepath.Join(path, name)}
			root.ChildFolders = append(root.ChildFolders, child)
		case 'l':
			continue
		}
	}
	respondJSON(w, http.StatusOK, root)
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
	log.Print(r.Method)
	switch r.Method {
	case http.MethodGet:
		l.fetch(w, r)
	case http.MethodPost:
		l.store(w, r)
	case http.MethodPut:
		fmt.Println(r.URL.Path)
		regex := regexp.MustCompile("/lessons/([0-9]+)")
		matches := regex.FindAllSubmatch([]byte(r.URL.Path), -1)
		fmt.Print(string(matches[0][1]))
		if len(matches) == 1 && len(matches[0]) == 2 {
			l.update(w, r, string(matches[0][1]))
		} else {
			respond(w, http.StatusNotFound)
		}
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
	ports := r.MultipartForm.Value["ports[]"]
	ports = append(ports, consolePort)
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
	db.Save(&lesson)
	if db.Error != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	u := user{}
	db.Where("id = ?", userID).First(&u)
	if db.Error != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	lessonDirectoryPath := filepath.Join(pathLessons, fmt.Sprintf("%d", lesson.ID))
	os.MkdirAll(lessonDirectoryPath, 0755)
	thumbnailHeaders := r.MultipartForm.File["thumbnail"]
	if len(thumbnailHeaders) == 1 {
		thumbnailHeader := thumbnailHeaders[0]
		filenameComponents := strings.Split(thumbnailHeader.Filename, ".")
		if 2 <= len(filenameComponents) {
			extension := filenameComponents[len(filenameComponents)-1]
			thumbnailPath := filepath.Join(pathLessonThumbnails, fmt.Sprintf("%d", lesson.ID), "thumbnail."+extension)
			os.MkdirAll(filepath.Dir(thumbnailPath), 0755)
			thumbnailFile, err := os.Create(filepath.Join(thumbnailPath))
			if err != nil {
				log.Println(err)
				respond(w, http.StatusInternalServerError)
				return
			}
			defer thumbnailFile.Close()
			thumbnail, err := thumbnailHeader.Open()
			if err != nil {
				log.Println(err)
				respond(w, http.StatusInternalServerError)
				return
			}
			io.Copy(thumbnailFile, thumbnail)
			lesson.ThumbnailPath = strings.TrimPrefix(thumbnailPath, cwd)
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
	_, err = d.exec(containername.String(), []string{"-d"}, "gotty", "-w", "-p", consolePort, "bash")
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	portsOutput, err := d.port(containername.String())
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	hostPorts := make(map[string]uint)
	for _, line := range strings.Split(string(portsOutput), "\n") {
		if line == "" {
			continue
		}
		components := strings.Split(line, " ")
		port := strings.Split(components[0], "/")[0]
		hostPort, err := strconv.ParseUint(strings.Split(components[2], ":")[1], 10, 64)
		hostPorts[port] = uint(hostPort)
		if err != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
		if port == consolePort {
			lesson.HostConsolePort = uint(hostPort)
		}
	}
	for _, port := range ports {
		portUint, err := strconv.ParseUint(port, 10, 64)
		if err != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
		lessonPort := lessonPort{
			Port:     uint(portUint),
			HostPort: hostPorts[port],
			LessonID: lesson.ID,
		}
		db.Save(&lessonPort)
		if db.Error != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
	}
	db.Save(&lesson)
	respondJSON(w, http.StatusOK, lesson)
}

func (l *lessonsHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	update(r, id, lesson{})
	if db.Error != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
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
