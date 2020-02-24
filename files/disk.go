package files

import (
	"io/ioutil"
	"os"
	"strings"
)

type DiskProvider struct{
	FileProvider
}

func (dp *DiskProvider) GetDirectory(path string) Directory {
	rp := strings.Join([]string{dp.Location,path}, "/")
	fileDir, err := ioutil.ReadDir(rp)
	if err != nil {
		_ = os.MkdirAll(path, 0644)
	}
	var fileList []FileInfo

	for _, file := range fileDir {
		info := FileInfo{
			IsDirectory: file.IsDir(),
			Name: file.Name(),
		}
		if !info.IsDirectory {
			split := strings.Split(file.Name(), ".")
			info.Extension = split[len(split) - 1]
		}
		fileList = append(fileList, info)
	}

	return Directory{
		Path: rp,
		Files: fileList,
	}
}

func (dp *DiskProvider) ViewFile(path string) string {
	return strings.Join([]string{dp.Location,path}, "/")
}

func (dp *DiskProvider) SaveFile(contents []byte, path string) bool {
	err := ioutil.WriteFile(path, contents, 0600)
	if err != nil {
		return false
	}
	return true
}