package main

import "crypto/sha512"

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
