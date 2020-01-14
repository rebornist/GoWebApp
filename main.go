package main

import (
	"fmt"
	"net/http"
	"webApp/router"
)

func main() {

	r := &router {
		make(map[string]map[string]http.HandlerFunc)
	}

	f.HandleFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "welcome!")
	})

	r.HandleFunc("GET", "/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "about")
	})

	r.HandleFunc("GET", "/users/:id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "retrieve user")
	})

	r.HandleFunc("GET", "/users/:user_id/addresses/:address_id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "retrieve user's address")
	})

	r.HandleFunc("POST", "/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "create user")
	})

	r.HandleFunc("POST", "/users/:user_id/addresses", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "create user's address")
	})

	// 8080 포트로 웹 서버 구동
	http.ListerAndServe(":8080", r)
}
