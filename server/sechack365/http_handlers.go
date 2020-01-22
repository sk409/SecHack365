package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/sk409/gofile"

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
	response := func(authenticated bool) map[string]bool {
		return map[string]bool{
			"authenticated": authenticated,
		}
	}
	sessionID, err := r.Cookie(cookieNameSessionID)
	if err != nil {
		respondJSON(w, http.StatusOK, response(false))
		return
	}
	_, err = sessionManager.Provider.Get(sessionID.Value)
	respondJSON(w, http.StatusOK, response(err == nil))
}

type downloadsHandler struct {
}

func (d *downloadsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		d.fetch(w, r)
		return
	case http.MethodPost:
		d.store(w, r)
		return
	}
	respond(w, http.StatusNotFound)
}

func (d *downloadsHandler) fetch(w http.ResponseWriter, r *http.Request) {
	downloads := []download{}
	err := fetch(r, &downloads)
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, downloads)
}

func (d *downloadsHandler) store(w http.ResponseWriter, r *http.Request) {
	userID := r.PostFormValue("userID")
	materialID := r.PostFormValue("materialID")
	if !notEmptyAll(userID, materialID) {
		respond(w, http.StatusBadRequest)
		return
	}
	userIDUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		respond(w, http.StatusBadRequest)
		return
	}
	materialIDUint, err := strconv.ParseUint(materialID, 10, 64)
	if err != nil {
		respond(w, http.StatusBadRequest)
		return
	}
	u := user{}
	db.Where("id = ?", userID).First(&u)
	if u.ID == 0 {
		respond(w, http.StatusBadRequest)
		return
	}
	mtl := material{}
	db.Where("id = ?", materialID).First(&mtl)
	if mtl.ID == 0 {
		respond(w, http.StatusBadRequest)
		return
	}
	mc := material{
		Title:         mtl.Title,
		Description:   mtl.Description,
		ThumbnailPath: strings.TrimPrefix(pathNoimage, cwd),
		Downloaded:    true,
		UserID:        uint(userIDUint),
		AuthorUserID:  &mtl.UserID,
	}
	db.Save(&mc)
	if db.Error != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	materialCloneThumbnailPath := filepath.Join(pathPublicImagesMaterials, fmt.Sprintf("%d", mc.ID), filepath.Base(mtl.ThumbnailPath))
	err = os.MkdirAll(filepath.Dir(materialCloneThumbnailPath), 0755)
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	err = gofile.CopyFileToFile(filepath.Join(cwd, mtl.ThumbnailPath), materialCloneThumbnailPath)
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	mc.ThumbnailPath = strings.TrimPrefix(materialCloneThumbnailPath, cwd)
	db.Save(&mc)
	lessonMaterials := []lessonMaterial{}
	db.Where("material_id = ?", mtl.ID).Find(&lessonMaterials)
	lessonIDs := make([]uint, len(lessonMaterials))
	for index, lessonMaterial := range lessonMaterials {
		lessonIDs[index] = lessonMaterial.LessonID
	}
	lessons := []lesson{}
	db.Where(lessonIDs).Find(&lessons)
	for _, l := range lessons {
		lessonPorts := []lessonPort{}
		db.Where("lesson_id = ?", l.ID).Find(&lessonPorts)
		if db.Error != nil {
			respond(w, http.StatusNotFound)
			return
		}
		ports := make([]string, len(lessonPorts))
		for index, lessonPort := range lessonPorts {
			ports[index] = fmt.Sprintf("%d", lessonPort.Port)
		}
		lc := lesson{
			Title:         l.Title,
			Description:   l.Description,
			Book:          l.Book,
			ThumbnailPath: strings.TrimPrefix(pathNoimage, cwd),
			ConsolePort:   l.ConsolePort,
			Downloaded:    true,
			UserID:        uint(userIDUint),
		}
		db.Save(&lc)
		if db.Error != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
		lessonCloneThumbnailPath := filepath.Join(pathPublicImagesLessons, fmt.Sprintf("%d", lc.ID), filepath.Base(l.ThumbnailPath))
		err = os.MkdirAll(filepath.Dir(lessonCloneThumbnailPath), 0755)
		if err != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
		err = gofile.CopyFileToFile(filepath.Join(cwd, l.ThumbnailPath), lessonCloneThumbnailPath)
		if err != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
		lc.ThumbnailPath = strings.TrimPrefix(lessonCloneThumbnailPath, cwd)
		d := docker{}
		imagename, err := uuid.NewUUID()
		if err != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
		err = d.commit(fmt.Sprintf("%s", l.DockerContainerID), imagename.String())
		if err != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
		consolePortString := fmt.Sprintf("%d", l.ConsolePort)
		dc, err := initDockerContainer(imagename.String(), consolePortString, ports...)
		if err != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
		lc.DockerContainerID = dc.id
		lc.HostConsolePort = dc.hostPorts[consolePortString]
		db.Save(&lc)
		for _, port := range ports {
			portUint, err := strconv.ParseUint(port, 10, 64)
			if err != nil {
				respond(w, http.StatusInternalServerError)
				return
			}
			lpc := lessonPort{
				Port:     uint(portUint),
				HostPort: dc.hostPorts[port],
				LessonID: lc.ID,
			}
			db.Save(&lpc)
			if db.Error != nil {
				respond(w, http.StatusInternalServerError)
				return
			}
		}
		lm := lessonMaterial{
			LessonID:   lc.ID,
			MaterialID: mc.ID,
		}
		db.Save(&lm)
		if db.Error != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
	}
	download := download{
		UserID:             uint(userIDUint),
		OriginalMaterialID: uint(materialIDUint),
		CopiedMaterialID:   mc.ID,
	}
	db.Save(&download)
	if db.Error != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, download)
}

type filesHandler struct {
}

func (f *filesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		f.fetch(w, r)
	case http.MethodPut:
		f.update(w, r)
	default:
		respond(w, http.StatusNotFound)
	}
}

func (f *filesHandler) fetch(w http.ResponseWriter, r *http.Request) {
	lessonID := r.URL.Query().Get("lessonID")
	path := r.URL.Query().Get("path")
	if !notEmptyAll(lessonID, path) {
		respond(w, http.StatusBadRequest)
		return
	}
	l := lesson{}
	db.Where("id = ?", lessonID).First(&l)
	if db.Error != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	d := docker{}
	text, err := d.exec(l.DockerContainerID, []string{}, "cat", path)
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	file := file{
		Dir:  filepath.Dir(path),
		Name: filepath.Base(path),
		Path: path,
		Text: string(text),
	}
	respondJSON(w, http.StatusOK, file)
}

func (f *filesHandler) update(w http.ResponseWriter, r *http.Request) {
	lessonID := r.PostFormValue("lessonID")
	path := r.PostFormValue("path")
	text := r.PostFormValue("text")
	if !notEmptyAll(lessonID, path) {
		respond(w, http.StatusBadRequest)
		return
	}
	l := lesson{}
	db.Where("id = ?", lessonID).First(&l)
	filename, err := uuid.NewUUID()
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	file, err := os.Create(filename.String())
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	defer file.Close()
	_, err = file.Write([]byte(text))
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	d := docker{}
	src := filepath.Join(cwd, filename.String())
	d.sendFile(l.DockerContainerID, src, path)
	os.Remove(filename.String())
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

type followsHandler struct {
}

func (f *followsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	base := "/follows/"
	switch r.Method {
	case http.MethodGet:
		f.fetch(w, r)
		return
	case http.MethodPost:
		f.store(w, r)
		return
	case http.MethodDelete:
		id, ok := routeWithID(r, base)
		if ok {
			f.destory(w, r, id)
			return
		}
	}
	respond(w, http.StatusNotFound)
}

func (f *followsHandler) fetch(w http.ResponseWriter, r *http.Request) {
	follows := []follow{}
	err := fetch(r, &follows)
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, follows)
}

func (f *followsHandler) store(w http.ResponseWriter, r *http.Request) {
	followingUserID := r.PostFormValue("followingUserID")
	followedUserID := r.PostFormValue("followedUserID")
	if !notEmptyAll(followingUserID, followedUserID) {
		respond(w, http.StatusBadRequest)
		return
	}
	followingUserIDUint, err := strconv.ParseUint(followingUserID, 10, 64)
	if err != nil {
		respond(w, http.StatusBadRequest)
		return
	}
	followedUserIDUint, err := strconv.ParseUint(followedUserID, 10, 64)
	if err != nil {
		respond(w, http.StatusBadRequest)
		return
	}
	count := 0
	db.Model(&user{}).Where("id = ?", followingUserID).Count(&count)
	if count == 0 {
		respond(w, http.StatusBadRequest)
		return
	}
	db.Model(&user{}).Where("id = ?", followedUserID).Count(&count)
	if count == 0 {
		respond(w, http.StatusBadRequest)
		return
	}
	follow := follow{
		FollowingUserID: uint(followingUserIDUint),
		FollowedUserID:  uint(followedUserIDUint),
	}
	db.Save(&follow)
	if db.Error != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, follow)
}

func (f *followsHandler) destory(w http.ResponseWriter, r *http.Request, id string) {
	db.Delete(follow{}, "id = ?", id)
	if db.Error != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
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
	base := "/lessons/"
	switch r.Method {
	case http.MethodGet:
		l.fetch(w, r)
		return
	case http.MethodPost:
		l.store(w, r)
		return
	case http.MethodPut:
		id, ok := routeWithID(r, base)
		if ok {
			l.update(w, r, id)
			return
		}
	case http.MethodDelete:
		id, ok := routeWithID(r, base)
		if ok {
			l.destroy(w, r, id)
			return
		}
	}
	respond(w, http.StatusNotFound)
}

func (l *lessonsHandler) fetch(w http.ResponseWriter, r *http.Request) {
	lessons := []lesson{}
	err := fetch(r, &lessons)
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, lessons)
}

func (l *lessonsHandler) store(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1 << 62)
	titles := r.MultipartForm.Value["title"]
	descriptions := r.MultipartForm.Value["description"]
	images := r.MultipartForm.Value["os"]
	superLessons := r.MultipartForm.Value["superLesson"]
	consolePorts := r.MultipartForm.Value["consolePort"]
	userIDs := r.MultipartForm.Value["userID"]
	if !notEmptyAll(titles, descriptions, consolePorts, userIDs) {
		respond(w, http.StatusBadRequest)
		return
	}
	if emptyAll(images, superLessons) {
		respond(w, http.StatusBadRequest)
		return
	}
	title := titles[0]
	description := descriptions[0]
	consolePort := consolePorts[0]
	userID := userIDs[0]
	ports := r.MultipartForm.Value["ports[]"]
	ports = append(ports, consolePort)
	if !notEmptyAll(title, description, consolePort, userID) {
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
	lsn := lesson{
		Title:         title,
		Description:   description,
		Book:          "",
		ConsolePort:   uint(consolePortUint),
		ThumbnailPath: strings.TrimPrefix(pathNoimage, cwd),
		UserID:        uint(userIDUint),
	}
	db.Save(&lsn)
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
	lessonDirectoryPath := filepath.Join(pathLessons, fmt.Sprintf("%d", lsn.ID))
	os.MkdirAll(lessonDirectoryPath, 0755)
	thumbnailHeaders := r.MultipartForm.File["thumbnail"]
	if len(thumbnailHeaders) == 1 {
		thumbnailPath := filepath.Join(pathPublicImagesLessons, fmt.Sprintf("%d", lsn.ID), "thumbnail")
		path, err := saveFile(thumbnailPath, thumbnailHeaders[0])
		if err != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
		lsn.ThumbnailPath = strings.TrimPrefix(path, cwd)
	}
	imagename := ""
	if len(superLessons) == 1 {
		superLesson := lesson{}
		db.Where("title = ? AND user_id = ?", superLessons[0], userID).First(&superLesson)
		if db.Error != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
		d := docker{}
		uid, err := uuid.NewUUID()
		if err != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
		imagename = uid.String()
		err = d.commit(superLesson.DockerContainerID, imagename)
		if err != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
	} else {
		image := images[0]
		imagename, err = buildDockerImage(image, u.Name, lessonDirectoryPath)
		if err != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
	}
	dc, err := initDockerContainer(imagename, consolePort, ports...)
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	lsn.DockerContainerID = dc.id
	lsn.HostConsolePort = dc.hostPorts[consolePort]
	for _, port := range ports {
		portUint, err := strconv.ParseUint(port, 10, 64)
		if err != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
		lessonPort := lessonPort{
			Port:     uint(portUint),
			HostPort: dc.hostPorts[port],
			LessonID: lsn.ID,
		}
		db.Save(&lessonPort)
		if db.Error != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
	}
	db.Save(&lsn)
	respondJSON(w, http.StatusOK, lsn)
}

func (l *lessonsHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	lsn := lesson{}
	db.Where("id = ?", id).First(&lsn)
	if lsn.ID == 0 {
		respond(w, http.StatusBadRequest)
		return
	}
	update(r, id, lesson{})
	if db.Error != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
}

func (l *lessonsHandler) destroy(w http.ResponseWriter, r *http.Request, id string) {
	err := deleteLesson(id)
	if err != nil {
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
	sessionCookie, err := r.Cookie(cookieNameSessionID)
	if err != nil {
		respond(w, http.StatusBadRequest)
		return
	}
	sessionID := sessionCookie.Value
	_, err = sessionManager.Provider.Get(sessionID)
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	err = sessionManager.Provider.Stop(sessionID)
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, newCookie(cookieNameSessionID, "", -1))
}

type materialsHandler struct {
}

func (m *materialsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	base := "/materials"
	switch r.Method {
	case http.MethodGet:
		m.fetch(w, r)
		return
	case http.MethodPost:
		m.store(w, r)
		return
	case http.MethodDelete:
		id, ok := routeWithID(r, base)
		if ok {
			m.destroy(w, r, id)
			return
		}
	}
	respond(w, http.StatusNotFound)
}

func (m *materialsHandler) fetch(w http.ResponseWriter, r *http.Request) {
	materials := []material{}
	err := fetch(r, &materials)
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, materials)
}

func (m *materialsHandler) store(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1 << 62)
	titles := r.MultipartForm.Value["title"]
	descriptions := r.MultipartForm.Value["description"]
	userIDs := r.MultipartForm.Value["userID"]
	if !notEmptyAll(titles, descriptions, userIDs) {
		respond(w, http.StatusBadRequest)
		return
	}
	title := titles[0]
	description := descriptions[0]
	userID := userIDs[0]
	lessonIDs := r.MultipartForm.Value["lessonIDs[]"]
	if !notEmptyAll(title, description, userID) {
		respond(w, http.StatusBadRequest)
		return
	}
	userIDUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		respond(w, http.StatusBadRequest)
		return
	}
	material := material{
		Title:         title,
		Description:   description,
		ThumbnailPath: strings.TrimPrefix(pathNoimage, cwd),
		UserID:        uint(userIDUint),
	}
	db.Save(&material)
	if db.Error != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	thumbnailHeaders := r.MultipartForm.File["thumbnail"]
	if len(thumbnailHeaders) == 1 {
		thumbnailPath := filepath.Join(pathPublicImagesMaterials, fmt.Sprintf("%d", material.ID), "thumbnail")
		path, err := saveFile(thumbnailPath, thumbnailHeaders[0])
		if err != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
		material.ThumbnailPath = strings.TrimPrefix(path, cwd)
		db.Save(&material)
	}
	for _, lessonID := range lessonIDs {
		lessonIDUint, err := strconv.ParseUint(lessonID, 10, 64)
		if err != nil {
			continue
		}
		count := 0
		db.Model(&lesson{}).Where("id = ?", lessonID).Count(&count)
		if count == 0 {
			continue
		}
		lm := lessonMaterial{
			LessonID:   uint(lessonIDUint),
			MaterialID: material.ID,
		}
		db.Save(&lm)
	}
	respondJSON(w, http.StatusOK, material)
}

// func (m *materialsHandler) download(w http.ResponseWriter, r *http.Request) {
// 	userID := r.PostFormValue("userID")
// 	materialID := r.PostFormValue("materialID")
// 	if !notEmptyAll(userID, materialID) {
// 		respond(w, http.StatusBadRequest)
// 		return
// 	}
// 	userIDUint, err := strconv.ParseUint(userID, 10, 64)
// 	if err != nil {
// 		respond(w, http.StatusBadRequest)
// 		return
// 	}
// 	u := user{}
// 	db.Where("id = ?", userID).First(&u)
// 	if u.ID == 0 {
// 		respond(w, http.StatusBadRequest)
// 		return
// 	}

// }

func (m *materialsHandler) destroy(w http.ResponseWriter, r *http.Request, id string) {
	mtl := material{}
	db.Where("id = ?", id).First(&mtl)
	if mtl.ID == 0 {
		respond(w, http.StatusBadRequest)
		return
	}
	if mtl.Downloaded {
		lessonMaterials := []lessonMaterial{}
		db.Where("material_id = ?", id).Find(&lessonMaterials)
		if db.Error != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
		for _, lm := range lessonMaterials {
			err := deleteLesson(fmt.Sprintf("%d", lm.LessonID))
			if err != nil {
				respond(w, http.StatusInternalServerError)
				return
			}
		}
	}
	db.Unscoped().Delete(material{}, "id = ?", id)
	if db.Error != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
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
	u := user{
		Name:             username,
		Password:         string(hashedPassword),
		ProfileImagePath: strings.TrimPrefix(pathNoimage, cwd),
	}
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
	base := "/users/"
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case base:
			u.fetch(w, r)
			return
		case path.Join(base, "search"):
			u.search(w, r)
			return
		default:
			id, ok := routeWithID(r, base)
			if ok {
				switch r.URL.Path {
				case base + id + "/follow":
					u.follow(w, r, id, false)
					return
				case base + id + "/follower":
					u.follow(w, r, id, true)
					return
				}
			}
		}
	case http.MethodPut:
		id, ok := routeWithID(r, base)
		if ok {
			u.update(w, r, id)
		}
	}
	respond(w, http.StatusNotFound)
}

func (u *usersHandler) fetch(w http.ResponseWriter, r *http.Request) {
	users := []user{}
	err := fetch(r, &users)
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, users)
}

func (u *usersHandler) follow(w http.ResponseWriter, r *http.Request, id string, follower bool) {
	follows := []follow{}
	if follower {
		db.Where("followed_user_id = ?", id).Find(&follows)
	} else {
		db.Where("following_user_id = ?", id).Find(&follows)
	}
	if db.Error != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	userIDs := make([]uint, len(follows))
	for index, follow := range follows {
		if follower {
			userIDs[index] = follow.FollowingUserID
		} else {
			userIDs[index] = follow.FollowedUserID
		}
	}
	users := []user{}
	db.Where(userIDs).Find(&users)
	if db.Error != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, users)
}

func (u *usersHandler) search(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	if !notEmptyAll(keyword) {
		respond(w, http.StatusBadRequest)
		return
	}
	users := []user{}
	db.Where("name LIKE ?", "%"+keyword+"%").Find(&users)
	if db.Error != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, users)
}

func (u *usersHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	r.ParseMultipartForm(1 << 62)
	thumbnailHeaders := r.MultipartForm.File["profileImage"]
	if len(thumbnailHeaders) == 1 {
		thumbnailPath := filepath.Join(pathPublicImagesUsers, id, "profileImage")
		path, err := saveFile(thumbnailPath, thumbnailHeaders[0])
		if err != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
		user := user{}
		db.Where("id = ?", id).First(&user)
		user.ProfileImagePath = strings.TrimPrefix(path, cwd)
		db.Save(&user)
	}
}
