package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/sk409/gotype"

	"github.com/sk409/goconst"
)

func fetch(r *http.Request, model interface{}) error {
	query := make(map[string]interface{})
	for key, values := range r.URL.Query() {
		query[key] = values[0]
	}
	db.Where(query).Find(model)
	return db.Error
}

func login(w http.ResponseWriter, u user) {
	session, err := sessionManager.Provider.Start()
	if err != nil {
		respond(w, http.StatusInternalServerError)
		return
	}
	session.Store(sessionKeyUser, u)
	http.SetCookie(w, newCookie(
		cookieNameSessionID,
		session.ID(),
		cookieMaxAge30Days,
	))
}

func public(data interface{}) (interface{}, error) {
	if gotype.IsSlice(data) {
		rv := reflect.ValueOf(data)
		s := make([]interface{}, rv.Len())
		for index := 0; index < rv.Len(); index++ {
			rvi := rv.Index(index).Interface()
			if f, ok := rvi.(facade); ok {
				p, err := f.public()
				if err != nil {
					return nil, err
				}
				s[index] = p
			} else {
				s[index] = rvi
			}
		}
		data = s
	}
	return data, nil
}

func respond(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

func respondJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	data, err := public(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	w.WriteHeader(statusCode)
	w.Write(jsonBytes)
}

func respondMessage(w http.ResponseWriter, statusCode int, message string) {
	respond(w, statusCode)
	w.Write([]byte(message))
}

func routeWithID(r *http.Request, base string) (string, bool) {
	if !strings.HasPrefix(base, "/") {
		base = "/" + base
	}
	if !strings.HasSuffix(base, "/") {
		base += "/"
	}
	regex := regexp.MustCompile(base + "([0-9]+)")
	matches := regex.FindAllSubmatch([]byte(r.URL.Path), -1)
	if len(matches) != 1 || len(matches[0]) != 2 {
		return "", false
	}
	return string(matches[0][1]), true
}

func saveFile(path string, header *multipart.FileHeader) (string, error) {
	filenameComponents := strings.Split(header.Filename, ".")
	extension := ""
	if 2 <= len(filenameComponents) {
		extension = filenameComponents[len(filenameComponents)-1]
	}
	path += "." + extension
	os.MkdirAll(filepath.Dir(path), 0755)
	file, err := os.Create(filepath.Join(path))
	if err != nil {
		return "", err
	}
	defer file.Close()
	multipartFile, err := header.Open()
	if err != nil {
		return "", err
	}
	io.Copy(file, multipartFile)
	return path, nil
}

func update(r *http.Request, id string, model interface{}) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	params := strings.Split(string(data), "&")
	query := make(map[string]interface{})
	for _, param := range params {
		keyvalue := strings.Split(param, "=")
		query[keyvalue[0]] = keyvalue[1]
	}
	db.Model(model).Where("id = ?", id).Updates(query)
	return db.Error
}
