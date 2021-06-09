package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func hashAndSalt(password string) string {
	bytesPasswd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytesPasswd, bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

func comparePassword(hashed string, compare []byte) bool {
	byteHash := []byte(hashed)
	err := bcrypt.CompareHashAndPassword(byteHash, compare)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
