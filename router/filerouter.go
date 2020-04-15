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
			filename, err := url.QueryUnescape(vars["file"])
			if err != nil {
				return &HTTPError{
					Message:    fmt.Sprintf("error determining filetype for %s\n", filename),
					StatusCode: http.StatusInternalServerError,
				}
			}
			ok, isDir, location := provider.ObjectInfo(filename)
			if !ok {
				return &HTTPError{
					Message:    fmt.Sprintf("error locating file %s\n", filename),
					StatusCode: http.StatusNotFound,
				}
			}

			if isDir {
				fileList := provider.GetDirectory(filename)

				data, err := json.Marshal(fileList)
				if err != nil {
					return &HTTPError{
						Message:    fmt.Sprintf("error fetching filelisting for %s\n", vars["file"]),
						StatusCode: http.StatusNotFound,
					}
				}
				w.Write(data)
				return nil
			}

			// If the file is local, attempt to use http.ServeContent for correct headers.
			if location == files.FILE_IS_LOCAL {
				rp := provider.FilePath(filename)
				if rp != "" {
					f, err := os.Open(rp)
					if err != nil {
						return &HTTPError{
							Message:    fmt.Sprintf("error opening file %s\n", rp),
							StatusCode: http.StatusInternalServerError,
						}
					}
					http.ServeContent(w, r, filename, time.Time{}, f)
				}
			}
			// If the file is remote, then delegate the writing to the response to the provider.
			// This isn't a great workaround, but avoids caching the whole file in mem or on disk.
			if location == files.FILE_IS_REMOTE {
				provider.RemoteFile(filename, w)
			}
			return nil
		}

		// File upload or directory creation.
		if r.Method == "POST" {
			xType := r.Header.Get("X-NAS-Type")
			// We only really care about this header of creating a directory.
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
				return nil
			}

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
			_, _  = w.Write([]byte("saved file"))
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
		_, _ = w.Write(data)
		return nil
	}
}
