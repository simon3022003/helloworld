package handlers

import (
	"net/http"
	"os"
)

// Hello returns "Hello World!" to the client
func Hello(w http.ResponseWriter, r *http.Request) {
	value, ok := os.LookupEnv("MESSAGE")
	w.WriteHeader(http.StatusOK)
	if ok {
		w.Write([]byte("Hello!" + value))
	} else {
		w.Write([]byte("Hello! world"))
	}
	return
}
