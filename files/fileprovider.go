package files

import (
	"fmt"
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

var Providers map[string]FileProvider

type FileContents struct {
	Content []byte
	IsUrl   bool
}

type FileProviderInterface interface {
	GetDirectory(path string) Directory
	ViewFile(path string, w io.Writer)
	SaveFile(contents []byte, path string) bool
	DetermineType(path string) string
}

func TranslateProvider(codename string, i *FileProviderInterface) {
	provider := Providers[codename]
	if codename == "disk" {
		*i = &DiskProvider{provider,}
		return
	}
	/*
	 * @TODO: It would be ideal if the authorization with Backblaze was done before
	 * actually needing to use it, ideally during the startup step.
	 */
	if codename == "backblaze" {
		bbProv := &BackblazeProvider{provider, provider.Config["bucket"], ""}

		err := bbProv.Authorize(provider.Config["appKeyId"], provider.Config["appId"])
		if err != nil {
			fmt.Println(err.Error())
		}
		*i = bbProv
		return
	}
	*i = FileProvider{}
}

/** DO NOT USE THESE DEFAULTS **/
func (f FileProvider) GetDirectory(path string) Directory {
	return Directory{}
}

func (f FileProvider) ViewFile(path string, w io.Writer) {
	return
}

func (f FileProvider) SaveFile(contents []byte, path string) bool {
	return false
}

func (f FileProvider) DetermineType(path string) string {
	return ""
}
