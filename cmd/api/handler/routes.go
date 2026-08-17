// filanem: cmd/api/handler/routes.go

package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *ApplicationDependencies) RegisterRoutes() http.Handler {
	// setup a router
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", a.home)
	router.HandlerFunc(http.MethodPost, "/api/jobs", a.CreateJobHandler)
	router.HandlerFunc(http.MethodGet, "/api/jobs", a.ShowAllJobsHandler)

	router.HandlerFunc(http.MethodPost, "/api/consumers", a.CreateConsumerHandler)
	router.HandlerFunc(http.MethodPatch, "/api/consumers/{id}", a.UpdateConsumerHandler)
	router.HandlerFunc(http.MethodGet, "/api/consumers", a.ShowAllConsumersHandler)

	router.HandlerFunc(http.MethodGet, "/api/api-keys", a.ShowAllAPIKeysHandler)
	router.HandlerFunc(http.MethodGet, "/api/api-keys/{id}", a.ShowAPIKeyHandler)
	router.HandlerFunc(http.MethodPost, "/api/api-keys", a.CreateAPIKeyHandler)

	recoverPanicMiddleware := a.recoverPanic(router)
	loggingMiddleware := a.loggingMiddleware(recoverPanicMiddleware)

	return loggingMiddleware

}
