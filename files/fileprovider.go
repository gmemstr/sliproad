package files

import (
	"io"
)

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
	Setup(args map[string]string) bool
	GetDirectory(path string) Directory
	ViewFile(path string) string
	RemoteFile(path string, writer io.Writer)
	SaveFile(file io.Reader, filename string, path string) bool
	ObjectInfo(path string) (string, string)
	CreateDirectory(path string) bool
	Delete(path string) bool
}

/** DO NOT USE THESE DEFAULTS **/
func (f FileProvider) Setup(args map[string]string) bool {
	return false
}

func (f FileProvider) GetDirectory(path string) Directory {
	return Directory{}
}

func (f FileProvider) ViewFile(path string) string {
	return ""
}

func (f FileProvider) RemoteFile(path string, writer io.Writer) {
	return
}

func (f FileProvider) SaveFile(file io.Reader, filename string, path string) bool {
	return false
}

func (f FileProvider) ObjectInfo(path string) (string, string) {
	return "", ""
}

func (f FileProvider) CreateDirectory(path string) bool {
	return false
}

func (f FileProvider) Delete(path string) bool {
	return false
}