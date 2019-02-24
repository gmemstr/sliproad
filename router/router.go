package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gmemstr/nas/common"
	"github.com/gmemstr/nas/files"
	"github.com/gorilla/mux"
)

type NewConfig struct {
	Name        string
	Host        string
	Email       string
	Description string
	Image       string
	PodcastURL  string
}

func Handle(handlers ...common.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		rc := &common.RouterContext{}
		for _, handler := range handlers {
			err := handler(rc, w, r)
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
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/web/static"))))

	// Paths that require specific handlers
	r.Handle("/", Handle(
		rootHandler(),
	)).Methods("GET")

	r.Handle("/files/", Handle(
		files.Listing("hot"),
	)).Methods("GET")
	r.Handle(`/files/{file:[a-zA-Z0-9=\-\/\s.,&]+}`, Handle(
		files.Listing("hot"),
	)).Methods("GET")
	r.Handle(`/file/{file:[a-zA-Z0-9=\-\/\s.,&]+}`, Handle(
		files.ViewFile("hot"),
	)).Methods("GET")

	r.Handle("/archive/", Handle(
		files.Listing("cold"),
	)).Methods("GET")
	r.Handle(`/archive/{file:[a-zA-Z0-9=\-\/\s.,&]+}`, Handle(
		files.Listing("cold"),
	)).Methods("GET")
	r.Handle(`/archived/{file:[a-zA-Z0-9=\-\/\s.,&]+}`, Handle(
		files.ViewFile("cold"),
	)).Methods("GET")


	return r
}

// Handles /.
func rootHandler() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {

		var file string
		switch r.URL.Path {
		case "/":
			w.Header().Set("Content-Type", "text/html")
			file = "assets/web/index.html"
		default:
			return &common.HTTPError{
				Message:    fmt.Sprintf("%s: Not Found", r.URL.Path),
				StatusCode: http.StatusNotFound,
			}
		}

		return common.ReadAndServeFile(file, w)
	}
}

func adminHandler() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		return common.ReadAndServeFile("assets/web/admin.html", w)
	}
}