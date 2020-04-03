package router

import (
	"encoding/json"
	"fmt"
	"github.com/gmemstr/nas/files"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
	"strings"
)

func HandleProvider() Handler {
	return func(context *Context, w http.ResponseWriter, r *http.Request) *HTTPError {
		vars := mux.Vars(r)
		providerCodename := vars["provider"]
		providerCodename = strings.Replace(providerCodename, "/", "", -1)
		provider := *files.Providers[providerCodename]

		if r.Method == "GET" {
			fileList := provider.GetDirectory("")
			if vars["file"] != "" {
				fileType := provider.DetermineType(vars["file"])
				if fileType == "" {
					return &HTTPError{
						Message:    fmt.Sprintf("error determining filetype for %s\n", vars["file"]),
						StatusCode: http.StatusInternalServerError,
					}
				}
				if fileType == "file" {
					provider.ViewFile(vars["file"], w)
					return nil
				}
				fileList = provider.GetDirectory(vars["file"])
			}
			data, err := json.Marshal(fileList)
			if err != nil {
				return &HTTPError{
					Message:    fmt.Sprintf("error fetching filelisting for %s\n", vars["file"]),
					StatusCode: http.StatusInternalServerError,
				}
			}
			w.Write(data)
		}

		if r.Method == "POST" {
			err := r.ParseMultipartForm(32 << 20)
			if err != nil {
				return &HTTPError{
					Message:    fmt.Sprintf("error parsing form for %s\n", vars["file"]),
					StatusCode: http.StatusInternalServerError,
				}
			}
			file, handler, err := r.FormFile("file")
			defer file.Close()

			success := provider.SaveFile(file, handler.Filename, vars["file"])
			if !success {
				return &HTTPError{
					Message:    fmt.Sprintf("error saving file %s\n", vars["file"]),
					StatusCode: http.StatusInternalServerError,
				}
			}
			w.Write([]byte("saved file"))
		}

		return nil
	}
}

func ListProviders() Handler {
	return func(context *Context, w http.ResponseWriter, r *http.Request) *HTTPError {
		var providers []string
		for v, _ := range files.ProviderConfig {
			providers = append(providers, v)
		}
		sort.Strings(providers)
		data, err := json.Marshal(providers)
		if err != nil {
			return &HTTPError{
				Message:    fmt.Sprintf("error provider listing"),
				StatusCode: http.StatusInternalServerError,
			}
		}
		w.Write(data)
		return nil
	}
}
