package router

import (
	"log"
	"net/http"
	"io/fs"

	"github.com/gorilla/mux"
)

type handler func(context *requestContext, w http.ResponseWriter, r *http.Request) *httpError

type httpError struct {
	Message    string
	StatusCode int
}

type requestContext struct{}

// Loop through passed functions and execute them, passing through the current
// requestContext, response writer and request reader.
func handle(handlers ...handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := &requestContext{}
		for _, handler := range handlers {
			err := handler(context, w, r)
			if err != nil {
				log.Printf("%v", err)
				w.Write([]byte(http.StatusText(err.StatusCode)))
				return
			}
		}
	})
}

// Init initializes the main router and all routes for the application.
func Init(sc fs.FS) *mux.Router {

	r := mux.NewRouter()

	// File & Provider API
	r.Handle("/api/providers", handle(
		requiresAuth(),
		listProviders(),
	)).Methods("GET")

	r.Handle(`/api/files/{provider:[a-zA-Z0-9]+\/*}`, handle(
		requiresAuth(),
		handleProvider(),
	)).Methods("GET", "POST", "DELETE")

	r.Handle(`/api/files/{provider:[a-zA-Z0-9]+}/{file:.+}`, handle(
		requiresAuth(),
		handleProvider(),
	)).Methods("GET", "POST", "DELETE")

	// Auth API & Endpoints
	r.Handle(`/api/auth/callback`, handle(
		callbackAuth(),
	)).Methods("GET", "POST")

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.FS(sc))))

	return r
}
