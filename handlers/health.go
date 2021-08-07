package handlers

import (
	"net/http"
	"time"
)

// Health returns 200 OK status code when called. Used for K8S liveness probe
func Health(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Second)
	w.WriteHeader(http.StatusOK)
	return
}
