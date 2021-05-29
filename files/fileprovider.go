package files

import (
	"io"
)

// FileIsRemote denotes whether file is a remote file.
const FileIsRemote = "remote"

// FileIsLocal denotes whether file is a local file.
const FileIsLocal = "local"

// FileProvider aggregates some very basic properties for authentication and
// provider decoding.
type FileProvider struct {
	Name           string            `yaml:"name"`
	Provider       string            `yaml:"provider"`
	Authentication string            `yaml:"authentication"`
	Location       string            `yaml:"path"`
	Config         map[string]string `yaml:"config"`
}

// Directory contains the path and a collection of FileInfos.
type Directory struct {
	Path  string
	Files []FileInfo
}

// FileInfo describes a single file or directory, doing it's best to
// figure out the extension as well.
type FileInfo struct {
	IsDirectory bool
	Name        string
	Extension   string
}

// FileProviderInterface provides some sane default functions.
type FileProviderInterface interface {
	Setup(args map[string]string) (ok bool)
	GetDirectory(path string) (directory Directory)
	SendFile(path string) (stream io.Reader, contenttype string, err error)
	SaveFile(file io.Reader, filename string, path string) (ok bool)
	ObjectInfo(path string) (exists bool, isDir bool, location string)
	CreateDirectory(path string) (ok bool)
	Delete(path string) (ok bool)
}

/** DO NOT USE THESE DEFAULTS **/

// Setup runs when the application starts up, and allows for things like authentication.
func (f FileProvider) Setup(args map[string]string) bool {
	return false
}

// GetDirectory fetches a directory's contents.
func (f FileProvider) GetDirectory(path string) Directory {
	return Directory{}
}

// RemoteFile will bypass http.ServeContent() and instead write directly to the response.
func (f FileProvider) SendFile(path string) (stream io.Reader, contenttype string, err error) {
	return
}

// SaveFile will save a file with the contents of the io.Reader at the path specified.
func (f FileProvider) SaveFile(file io.Reader, filename string, path string) bool {
	return false
}

// ObjectInfo will return the info for an object given a path to if the file exists and location.
// Should return whether the path exists, if the path is a directory, and if it lives on disk.
// (see constants defined: `FILE_IS_REMOTE` and `FILE_IS_LOCAL`)
func (f FileProvider) ObjectInfo(path string) (bool, bool, string) {
	return false, false, ""
}

// CreateDirectory will create a directory on services that support it.
func (f FileProvider) CreateDirectory(path string) bool {
	return false
}

// Delete simply deletes a file. This is expected to be a destructive action by default.
func (f FileProvider) Delete(path string) bool {
	return false
}
