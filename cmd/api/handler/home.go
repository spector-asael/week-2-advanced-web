// Filename: cmd/api/handler/home.go

package handler

import (
	"net/http"
)

func (a *ApplicationDependencies) home(w http.ResponseWriter, r *http.Request) {
	// prepare a response to send to the client
	message := "Welcome to the Advanced Web Development API"
	a.writeJSON(w, http.StatusOK, envelope{"message": message}, nil)
}
