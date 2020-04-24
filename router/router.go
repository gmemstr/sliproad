package router

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

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
func Init() *mux.Router {

	r := mux.NewRouter()

	// "Static" paths
	r.PathPrefix("/javascript/").Handler(http.StripPrefix("/javascript/", http.FileServer(http.Dir("assets/web/javascript"))))
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("assets/web/css"))))
	r.PathPrefix("/icons/").Handler(http.StripPrefix("/icons/", http.FileServer(http.Dir("assets/web/icons"))))

	// Paths that require specific handlers
	r.Handle("/", handle(
		requiresAuth(),
		rootHandler(),
	)).Methods("GET")

	// File & Provider API
	r.Handle("/api/providers", handle(
		requiresAuth(),
		listProviders(),
	)).Methods("GET")

	r.Handle(`/api/files/{provider:[a-zA-Z0-9]+\/*}`, handle(
		requiresAuth(),
		handleProvider(),
	)).Methods("GET", "POST")

	r.Handle(`/api/files/{provider:[a-zA-Z0-9]+}/{file:.+}`, handle(
		requiresAuth(),
		handleProvider(),
	)).Methods("GET", "POST", "DELETE")

	// Auth API & Endpoints
	r.Handle(`/api/auth/callback`, handle(
		callbackAuth(),
	)).Methods("GET", "POST")

	return r
}

// Handles serving index page.
func rootHandler() handler {
	return func(context *requestContext, w http.ResponseWriter, r *http.Request) *httpError {
		f, err := os.Open("assets/web/index.html")
		if err != nil {
			return &httpError{
				Message:    fmt.Sprintf("error serving index page from assets/web"),
				StatusCode: http.StatusInternalServerError,
			}
		}

		defer f.Close()
		stats, err := f.Stat()
		if err != nil {
			return &httpError{
				Message:    fmt.Sprintf("error serving index page from assets/web"),
				StatusCode: http.StatusInternalServerError,
			}
		} else {
			w.Header().Add("Content-Length", strconv.FormatInt(stats.Size(), 10))
		}

		_, err = io.Copy(w, f)
		if err != nil {
			return &httpError{
				Message:    fmt.Sprintf("error serving index page from assets/web"),
				StatusCode: http.StatusInternalServerError,
			}
		}
		return nil
	}
}
