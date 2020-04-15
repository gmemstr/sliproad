package files

import (
	"io"
)

const FILE_IS_REMOTE = "remote"
const FILE_IS_LOCAL = "local"

type FileProvider struct {
	Name           string            `yaml:"name"`
	Provider       string            `yaml:"provider"`
	Authentication string            `yaml:"authentication"`
	Location       string            `yaml:"path"`
	Config         map[string]string `yaml:"config"`
}

type Directory struct {
	Path  string
	Files []FileInfo
}

type FileInfo struct {
	IsDirectory bool
	Name        string
	Extension   string
}

type FileContents struct {
	Content []byte
	IsUrl   bool
}

type FileProviderInterface interface {
	// Called on initial startup of the application.
	Setup(args map[string]string) (ok bool)
	// Fetches the contents of a "directory".
	GetDirectory(path string) (directory Directory)
	// Builds a path to a file, for serving.
	FilePath(path string) (realpath string)
	// Fetch and pass along a remote file directly to the response writer.
	RemoteFile(path string, writer io.Writer)
	// Save a file from an io.Reader.
	SaveFile(file io.Reader, filename string, path string) (ok bool)
	// Fetch the info for an object given a path to if the file exists and location.
	// Should return whether the path exists, if the path is a directory, and if it lives on disk.
	// (see constants defined: `FILE_IS_REMOTE` and `FILE_IS_LOCAL`)
	ObjectInfo(path string) (exists bool, isDir bool, location string)
	// Create a directory if possible, returns the result.
	CreateDirectory(path string) (ok bool)
	// Delete a file or directory.
	Delete(path string) (ok bool)
}

/** DO NOT USE THESE DEFAULTS **/
func (f FileProvider) Setup(args map[string]string) bool {
	return false
}

func (f FileProvider) GetDirectory(path string) Directory {
	return Directory{}
}

func (f FileProvider) FilePath(path string) string {
	return ""
}

func (f FileProvider) RemoteFile(path string, writer io.Writer) {
	return
}

func (f FileProvider) SaveFile(file io.Reader, filename string, path string) bool {
	return false
}

func (f FileProvider) ObjectInfo(path string) (bool, bool, string) {
	return false, false, ""
}

func (f FileProvider) CreateDirectory(path string) bool {
	return false
}

func (f FileProvider) Delete(path string) bool {
	return false
}