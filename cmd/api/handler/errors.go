// Filename: cmd/api/errors.go
package handler

import (
	"errors"
	"net/http"
)

var ErrRecordNotFound = errors.New("record not found")

// log an error message
func (a *ApplicationDependencies) logError(r *http.Request, err error) {

	method := r.Method
	uri := r.URL.RequestURI()
	a.Logger.Error(err.Error(), "method", method, "uri", uri)

}

// send an error response in JSON
func (a *ApplicationDependencies) errorResponseJSON(w http.ResponseWriter,
	r *http.Request,
	status int,
	message any) {

	errorData := envelope{"error": message}
	err := a.writeJSON(w, status, errorData, nil)
	if err != nil {
		a.logError(r, err)
		w.WriteHeader(500)
	}
}

func (a *ApplicationDependencies) errorResponse(w http.ResponseWriter, r *http.Request, errorStatus int, message string) {
	errorData := envelope{"error": message}
	err := a.writeJSON(w, errorStatus, errorData, nil)
	if err != nil {
		a.logError(r, err)
		w.WriteHeader(500)
	}
}

// send an error response if our server messes up
func (a *ApplicationDependencies) serverErrorResponse(w http.ResponseWriter,
	r *http.Request,
	err error) {

	// first thing is to log error message
	a.logError(r, err)
	// prepare a response to send to the client
	message := "the server encountered a problem and could not process your request"
	a.errorResponseJSON(w, r, http.StatusInternalServerError, message)
}

func (app *ApplicationDependencies) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (a *ApplicationDependencies) badRequestResponse(w http.ResponseWriter,
	r *http.Request,
	err error) {

	a.errorResponseJSON(w, r, http.StatusBadRequest, err.Error())
}

func (a *ApplicationDependencies) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	a.errorResponseJSON(w, r, http.StatusUnprocessableEntity, errors)
}
