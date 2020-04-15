package router

import (
	"encoding/json"
	"fmt"
	"github.com/gmemstr/nas/files"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

func HandleProvider() Handler {
	return func(context *Context, w http.ResponseWriter, r *http.Request) *HTTPError {
		vars := mux.Vars(r)
		providerCodename := vars["provider"]
		providerCodename = strings.Replace(providerCodename, "/", "", -1)
		provider := *files.Providers[providerCodename]

		// Directory listing or serve file.
		if r.Method == "GET" {
			fileList := provider.GetDirectory("")
			if vars["file"] != "" {
				filename, err := url.QueryUnescape(vars["file"])
				if err != nil {
					return &HTTPError{
						Message:    fmt.Sprintf("error determining filetype for %s\n", filename),
						StatusCode: http.StatusInternalServerError,
					}
				}
				fileType, location := provider.ObjectInfo(filename)

				if fileType == "" {
					return &HTTPError{
						Message:    fmt.Sprintf("error determining filetype for %s\n", filename),
						StatusCode: http.StatusInternalServerError,
					}
				}
				if fileType == "file" {
					if location == "local" {
						rp := provider.ViewFile(filename)
						if rp != "" {
							f, _ := os.Open(rp)
							http.ServeContent(w, r, filename, time.Time{}, f)
						}
					}
					if location == "remote" {
						provider.RemoteFile(filename, w)
					}
					return nil
				}
				fileList = provider.GetDirectory(filename)
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

		// File upload or directory creation.
		if r.Method == "POST" {
			xType := r.Header.Get("X-NAS-Type")

			if xType == "file" || xType == ""{
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
			if xType == "directory" {
				dirname := vars["file"]
				success := provider.CreateDirectory(dirname)
				if !success {
					return &HTTPError{
						Message:    fmt.Sprintf("error creating directory %s\n", dirname),
						StatusCode: http.StatusInternalServerError,
					}
				}
				_, _ = w.Write([]byte("created directory"))
			}
		}

		// Delete file.
		if r.Method == "DELETE" {
			path := vars["file"]
			success := provider.Delete(path)
			if !success {
				return &HTTPError{
					Message:    fmt.Sprintf("error deleting %s\n", path),
					StatusCode: http.StatusInternalServerError,
				}
			}
			_, _ = w.Write([]byte("deleted"))
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
