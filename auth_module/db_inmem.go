package main

import (
	"math/rand"
)

type InmemUser struct {
	User
	uid string
}

type Inmem struct {
	userList []InmemUser
	session  map[string]string
}

func (inm *Inmem) Login(username, password string) (bool, string, error) {
	for _, v := range inm.userList {
		if v.username == username && v.password == password {
			return true, v.uid, nil
		}
	}
	return false, "", nil
}

func (inm *Inmem) CreateSession(uid string) (string, error) {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomStringLen := 16
	b := make([]byte, randomStringLen)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	randomString := string(b)
	inm.session[uid] = randomString
	return randomString, nil
}
