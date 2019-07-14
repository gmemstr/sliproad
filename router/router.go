package router

import (
	"database/sql"
	"fmt"
	"github.com/gmemstr/nas/auth"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"

	"github.com/gmemstr/nas/common"
	"github.com/gmemstr/nas/files"
	"github.com/gmemstr/nas/system"
	"github.com/gorilla/mux"
)

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
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("assets/web/static"))))

	// Paths that require specific handlers
	r.Handle("/", Handle(
		auth.RequireAuthorization(1),
		rootHandler(),
	)).Methods("GET")

	r.Handle(`/login`, Handle(
		loginHandler(),
	)).Methods("POST", "GET")

	r.Handle("/api/diskusage", Handle(
		auth.RequireAuthorization(1),
		system.DiskUsages(),
	)).Methods("GET")

	r.Handle(`/api/file/{file:[a-zA-Z0-9=\-\/\s.,&_+]+}`, Handle(
		auth.RequireAuthorization(1),
		files.ViewFile(),
	)).Methods("GET")
	r.Handle("/api/upload", Handle(
		files.UploadFile(),
	)).Methods("POST")

	r.Handle("/api/{tier:(?:hot|cold)}/", Handle(
		auth.RequireAuthorization(1),
		files.Listing(),
	)).Methods("GET")
	r.Handle(`/api/{tier:^(?:hot|cold)$}/{file:[a-zA-Z0-9=\-\/\s.,&_+]+}`, Handle(
		auth.RequireAuthorization(1),
		files.Listing(),
	)).Methods("GET")

	return r
}


func loginHandler() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "text/html")
			file := "assets/web/index.html"

			return common.ReadAndServeFile(file, w)
		}
		db, err := sql.Open("sqlite3", "assets/config/users.db")

		if err != nil {
			return &common.HTTPError{
				Message:    fmt.Sprintf("error in reading user database: %v", err),
				StatusCode: http.StatusInternalServerError,
			}
		}

		statement, err := db.Prepare("SELECT * FROM users WHERE username=?")

		if _, err := auth.DecryptCookie(r); err == nil {
			http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
			return nil
		}

		err = r.ParseForm()
		if err != nil {
			return &common.HTTPError{
				Message:    fmt.Sprintf("error in parsing form: %v", err),
				StatusCode: http.StatusBadRequest,
			}
		}

		username := r.Form.Get("username")
		password := r.Form.Get("password")
		rows, err := statement.Query(username)

		if username == "" || password == "" || err != nil {
			return &common.HTTPError{
				Message:    "username or password is invalid",
				StatusCode: http.StatusBadRequest,
			}
		}
		var id int
		var dbun string
		var dbhsh string
		var dbrn string
		var dbem string
		var dbperm int
		for rows.Next() {
			err := rows.Scan(&id, &dbun, &dbhsh, &dbrn, &dbem, &dbperm)
			if err != nil {
				return &common.HTTPError{
					Message:    fmt.Sprintf("error in decoding sql data", err),
					StatusCode: http.StatusBadRequest,
				}
			}

		}
		// Create a cookie here because the credentials are correct
		if bcrypt.CompareHashAndPassword([]byte(dbhsh), []byte(password)) == nil {
			c, err := auth.CreateSession(&common.User{
				Username: username,
			})
			if err != nil {
				return &common.HTTPError{
					Message:    err.Error(),
					StatusCode: http.StatusInternalServerError,
				}
			}

			// r.AddCookie(c)
			w.Header().Add("Set-Cookie", c.String())
			// And now redirect the user to admin page
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			db.Close()
			return nil
		}

		return &common.HTTPError{
			Message:    "Invalid credentials!",
			StatusCode: http.StatusUnauthorized,
		}
	}
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