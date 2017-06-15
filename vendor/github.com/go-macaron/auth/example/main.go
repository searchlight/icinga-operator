package main

import (
	"github.com/go-macaron/auth"
	"gopkg.in/macaron.v1"
)

func main() {
	m := macaron.Classic()
	// authenticate every request
	m.Use(auth.BasicFunc(func(username, password string) bool {
		return username == "admin" && password == "guessme"
	}))
	m.Get("/", func() string {
		return "Hello World!"
	})
	m.Run()
}
