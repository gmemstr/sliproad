package files

import (
	"encoding/json"
	"fmt"
	"github.com/gmemstr/nas/auth"
	"github.com/gmemstr/nas/common"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	ColdStorage string
	HotStorage  string
}

type Directory struct {
	Path         string
	Files        []FileInfo
	Previous     string
	Prefix       string
	SinglePrefix string
}

type FileInfo struct {
	IsDirectory bool
	Name        string
}

func GetUserDirectory(r *http.Request, tier string) (string, string, string) {
	usr, err := auth.DecryptCookie(r)
	if err != nil {
		return "", "", ""
	}

	username := usr.Username

	d, err := ioutil.ReadFile("assets/config/config.json")
	if err != nil {
		panic(err)
	}

	var config Config
	err = json.Unmarshal(d, &config)
	if err != nil {
		panic(err)
	}

	// Default to hot storage
	storage := config.HotStorage + username
	prefix := "files"
	singleprefix := "file"
	if tier == "cold" {
		storage = config.ColdStorage + username
		prefix = "archive"
		singleprefix = "archived"
	}

	return storage, prefix, singleprefix
}

// Lists out directory using template.
func Listing() common.Handler {

	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		vars := mux.Vars(r)
		id := vars["file"]
		tier := vars["tier"]
		storage, prefix, singleprefix := GetUserDirectory(r, tier)
		if storage == "" && prefix == "" && singleprefix == "" {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return &common.HTTPError{
				Message:    "Unauthorized, or unable to find cookie",
				StatusCode: http.StatusTemporaryRedirect,
			}
		}
		path := storage
		if id != "" {
			path = storage + id
		}

		fileDir, err := ioutil.ReadDir(path)
		if err != nil && path == "" {
			fmt.Println(path)
			_ = os.MkdirAll(path, 0644)
		}
		var fileList []FileInfo

		for _, file := range fileDir {
			info := FileInfo{
				IsDirectory: file.IsDir(),
				Name: file.Name(),
			}
			fileList = append(fileList, info)
		}
		path = strings.Replace(path, storage, "", -1)

		// Figure out what our previous location was.
		previous := strings.Split(path, "/")
		previous = previous[:len(previous)-1]
		previousPath := strings.Join(previous, "/")

		directory := Directory{
			Path:         path,
			Files:        fileList,
			Previous:     previousPath,
			Prefix:       prefix,
			SinglePrefix: singleprefix,
		}

		resultJson, err := json.Marshal(directory);
		w.Write(resultJson);
		return nil;
	}

}

// Lists out directory using template.
func ViewFile() common.Handler {

	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		vars := mux.Vars(r)
		id := vars["file"]
		tier := vars["tier"]

		d, err := ioutil.ReadFile("assets/config/config.json")
		if err != nil {
			panic(err)
		}

		var config Config;
		err = json.Unmarshal(d, &config)
		if err != nil {
			panic(err)
		}
		// Default to hot storage
		storage, _, _ := GetUserDirectory(r, tier)
		path := storage + id

		common.ReadAndServeFile(path, w)
		return nil
	}

}

func UploadFile() common.Handler {

	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		d, err := ioutil.ReadFile("assets/config/config.json")
		var config Config;
		err = json.Unmarshal(d, &config)
		if err != nil {
			panic(err)
		}

		err = r.ParseMultipartForm(32 << 20)
		path := strings.Join(r.Form["path"], "")

		// Default to hot storage
		storage := config.HotStorage

		file, handler, err := r.FormFile("file")
		if err != nil {
			fmt.Println(err)
			return nil
		}
		defer file.Close()

		f, err := os.OpenFile(storage+"/"+path+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)

		if err != nil {
			panic(err)
		}
		defer f.Close()
		io.Copy(f, file)
		return nil
	}

}
