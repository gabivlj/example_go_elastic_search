package main

import (
	"math/rand"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Message struct {
	User          string    `json:"user"`
	Text          string    `json:"textContent"`
	Date          time.Time `json:"date"`
	IndexPosition int64     `json:"indexPosition"`
	Users         []string  `json:"users"`
}

func RandomString() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandomStrings(len int) []string {
	strings := make([]string, 0)
	for i := 0; i < len; i++ {
		strings = append(strings, RandomString())
	}
	return strings
}

func ChooseString(s []string) string {
	rand.Seed(time.Now().UnixNano())
	return s[rand.Intn(len(s))]
}
