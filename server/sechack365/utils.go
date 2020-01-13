package main

import (
	"crypto/sha512"
	"reflect"
)

func encrypt(data []byte) []byte {
	hash := sha512.Sum512(data)
	return hash[:]
}

func notEmptyAll(vars ...string) bool {
	for _, v := range vars {
		if v == "" {
			return false
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
