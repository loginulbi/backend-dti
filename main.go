package main

import (
	"net/http"

	"login-service/route"
)

func main() {
	http.HandleFunc("/", route.URL)
	http.ListenAndServe(":8080", nil)
}
