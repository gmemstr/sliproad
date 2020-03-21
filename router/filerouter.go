package router

import (
	"encoding/json"
	"github.com/gmemstr/nas/common"
	"github.com/gmemstr/nas/files"
	"github.com/gorilla/mux"
	"net/http"
)

func HandleProvider() common.Handler {

	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		vars := mux.Vars(r)
		if r.Method == "GET" {
			providerCodename := vars["provider"]
			provider := *files.Providers[providerCodename]

			fileList := provider.GetDirectory("")
			if vars["file"] != "" {
				fileType := provider.DetermineType(vars["file"])
				if fileType == "" {
					w.Write([]byte("file not found"))
					return nil
				}
				if fileType == "file" {
					provider.ViewFile(vars["file"], w)
					return nil
				}
				fileList = provider.GetDirectory(vars["file"])
			}
			data, err := json.Marshal(fileList)
			if err != nil {
				w.Write([]byte("An error occurred"))
				return nil
			}
			w.Write(data)
		}

		return nil
	}
}

func ListProviders() common.Handler {

	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		return nil
	}
}
