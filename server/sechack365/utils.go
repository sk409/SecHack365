package main

import (
	"crypto/sha512"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/sk409/gotype"
)

func authUser(r *http.Request) (*user, error) {
	cookie, err := r.Cookie(cookieNameSessionID)
	if err != nil {
		return nil, err
	}
	session, err := sessionManager.Provider.Get(cookie.Value)
	if err != nil {
		return nil, err
	}
	u := user{}
	err = session.Object(sessionKeyUser, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func buildDockerImage(image, username, directoryPath string) (string, error) {
	df := newDockerfile(image, username)
	dockerfilePath := filepath.Join(directoryPath, "Dockerfile")
	err := df.write(dockerfilePath)
	if err != nil {
		return "", err
	}
	imagename, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	d := docker{}
	err = d.buildImage(imagename.String(), filepath.Dir(dockerfilePath))
	if err != nil {
		return "", err
	}
	return imagename.String(), nil
}

func deleteLesson(id string) error {
	lsn := lesson{}
	db.Where("id = ?", id).First(&lsn)
	if lsn.ID == 0 {
		return errBadRequest
	}
	db.Unscoped().Delete(lesson{}, "id = ?", id)
	if db.Error != nil {
		return db.Error
	}
	d := docker{}
	_, err := d.kill(lsn.DockerContainerID)
	if err != nil {
		return err
	}
	_, err = d.removeContainer(lsn.DockerContainerID)
	if err != nil {
		return err
	}
	return nil
}

func emptyAll(values ...interface{}) bool {
	for _, value := range values {
		if gotype.IsString(value) {
			if value.(string) != "" {
				return false
			}
		} else if gotype.IsSlice(value) {
			rv := reflect.ValueOf(value)
			if rv.Len() != 0 {
				return false
			}
		}
	}
	return true
}

func encrypt(data []byte) []byte {
	hash := sha512.Sum512(data)
	return hash[:]
}

func initDockerContainer(imagename, consolePort string, ports ...string) (*dockerContainer, error) {
	d := docker{}
	containername, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	containerID, err := d.runContainer(containername.String(), imagename, ports...)
	if err != nil {
		return nil, err
	}
	_, err = d.exec(containername.String(), []string{"-d"}, "gotty", "-w", "-p", consolePort, "bash")
	if err != nil {
		return nil, err
	}
	// _, err = d.exec(containername.String(), []string{"-d"}, "find /var/lib/mysql -type f -exec touch {} \\;")
	// if err != nil {
	// 	return nil, err
	// }
	portsOutput, err := d.port(containername.String())
	if err != nil {
		return nil, err
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
			return nil, err
		}
	}
	dc := dockerContainer{
		id:        string(containerID),
		name:      containername.String(),
		hostPorts: hostPorts,
	}
	return &dc, nil
}

func notEmptyAll(vars ...interface{}) bool {
	for _, v := range vars {
		if gotype.IsString(v) {
			if v.(string) == "" {
				return false
			}
		} else if gotype.IsSlice(v) {
			rv := reflect.ValueOf(v)
			if rv.Len() == 0 {
				return false
			}
		}
	}
	return true
}

func structToMap(data interface{}) map[string]interface{} {
	rt := reflect.TypeOf(data)
	rv := reflect.ValueOf(data)
	m := make(map[string]interface{})
	for i := 0; i < rt.NumField(); i++ {
		ct := rt.Field(i)
		cv := rv.Field(i)
		m[ct.Name] = cv.Interface()
	}
	return m
}
