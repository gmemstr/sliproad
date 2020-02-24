package files

import "fmt"

type FileProvider struct {
	Name string `yaml:"name"`
	Authentication string `yaml:"authentication"`
	Location string `yaml:"path"`
	Config map[string]string `yaml:"config"`
}

type Directory struct {
	Path         string
	Files        []FileInfo
}

type FileInfo struct {
	IsDirectory bool
	Name        string
	Extension   string
}

var Providers map[string]FileProvider

type FileProviderInterface interface {
	GetDirectory(path string) Directory
	ViewFile(path string) string
	SaveFile(contents []byte, path string) bool
}

func TranslateProvider(codename string, i *FileProviderInterface) {
	provider := Providers[codename]
	if codename == "disk" {
		*i = &DiskProvider{provider,}
		return
	}
	if codename == "backblaze" {
		bbProv := &BackblazeProvider{provider, provider.Config["bucket"]}

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

func (f FileProvider) ViewFile(path string) string {
	return ""
}

func (f FileProvider) SaveFile(contents []byte, path string) bool {
	return false
}