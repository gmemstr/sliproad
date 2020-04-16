package router

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Handler func(context *Context, w http.ResponseWriter, r *http.Request) *HTTPError

type HTTPError struct {
	Message string
	StatusCode int
}

// Context contains any information to be shared with middlewares.
type Context struct {}

func Handle(handlers ...Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := &Context{}
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


// Actual router, define endpoints here.
func Init() *mux.Router {

	r := mux.NewRouter()

	// "Static" paths
	r.PathPrefix("/javascript/").Handler(http.StripPrefix("/javascript/", http.FileServer(http.Dir("assets/web/javascript"))))
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("assets/web/css"))))
	r.PathPrefix("/icons/").Handler(http.StripPrefix("/icons/", http.FileServer(http.Dir("assets/web/icons"))))

	// Paths that require specific handlers
	r.Handle("/", Handle(
		requiresAuth(),
		rootHandler(),
	)).Methods("GET")

	// File & Provider API
	r.Handle("/api/providers", Handle(
		requiresAuth(),
		ListProviders(),
	)).Methods("GET")
	
	r.Handle(`/api/files/{provider:[a-zA-Z0-9]+\/*}`, Handle(
		requiresAuth(),
		HandleProvider(),
	)).Methods("GET", "POST")

	r.Handle(`/api/files/{provider:[a-zA-Z0-9]+}/{file:.+}`, Handle(
		requiresAuth(),
		HandleProvider(),
	)).Methods("GET", "POST", "DELETE")

	// Auth API & Endpoints
	r.Handle(`/api/auth/callback`, Handle(
		callbackAuth(),
	)).Methods("GET", "POST")

	return r
}

// Handles serving index page.
func rootHandler() Handler {
	return func(context *Context, w http.ResponseWriter, r *http.Request) *HTTPError {
		f, err := os.Open("assets/web/index.html")
		if err != nil {
			return &HTTPError{
				Message:    fmt.Sprintf("error serving index page from assets/web"),
				StatusCode: http.StatusInternalServerError,
			}
		}

		defer f.Close()
		stats, err := f.Stat()
		if err != nil {
			return &HTTPError{
				Message:    fmt.Sprintf("error serving index page from assets/web"),
				StatusCode: http.StatusInternalServerError,
			}
		} else {
			w.Header().Add("Content-Length", strconv.FormatInt(stats.Size(), 10))
		}

		_, err = io.Copy(w, f)
		if err != nil {
			return &HTTPError{
				Message:    fmt.Sprintf("error serving index page from assets/web"),
				StatusCode: http.StatusInternalServerError,
			}
		}
		return nil
	}
}