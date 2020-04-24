package router

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gmemstr/nas/authentication"
)

// AuthEnabled is a global variable that determines whether we were
// able to set up an authentication method (e.g Keycloak).
var AuthEnabled bool = true

func requiresAuth() handler {
	return func(context *requestContext, w http.ResponseWriter, r *http.Request) *httpError {
		if !AuthEnabled {
			return nil
		}
		cookie, err := r.Cookie("NAS-SESSION")
		if err != nil || !authentication.HasAuth(cookie.Value) {
			if err != nil {
				fmt.Println("Error", err.Error())
			}
			http.Redirect(w, r, authentication.GetLoginLink(), 307)
			return &httpError{
				Message:    "Unauthorized! Redirecting to /login",
				StatusCode: http.StatusTemporaryRedirect,
			}
		}
		return nil
	}
}

func callbackAuth() handler {
	return func(context *requestContext, w http.ResponseWriter, r *http.Request) *httpError {
		// Translate callback GET to POST to set cookie, then redirect.
		if r.Method == "GET" {
			javascript := `
<script>fetch("/api/auth/callback", {method:"POST", body: window.location.hash.split("&")[1].split("=")[1]}).then((r) => window.location.href = "/")</script>`
			w.Write([]byte(javascript))
			return nil
		}
		token, _ := ioutil.ReadAll(r.Body)

		// Set as HttpOnly cookie to mitigate XSS risk.
		jwtCookie := http.Cookie{Name: "NAS-SESSION",
			Value:    string(token),
			HttpOnly: true,
			Path:     "/",
		}

		http.SetCookie(w, &jwtCookie)

		return nil
	}
}
