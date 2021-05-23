package files

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type DiskProvider struct{
	FileProvider
}

func (dp *DiskProvider) Setup(args map[string]string) bool {
	return true
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

func (dp *DiskProvider) SendFile(path string, writer io.Writer) (stream io.Reader, contenttype string, err error) {
	rp := strings.Join([]string{dp.Location,path}, "/")
	f, err := os.Open(rp)
	if err != nil {
		return stream, contenttype, err
	}

	contenttype = mime.TypeByExtension(filepath.Ext(rp))

	if contenttype == "" {
		var buf [512]byte
		n, _ := io.ReadFull(f, buf[:])
		contenttype = http.DetectContentType(buf[:n])
	}

	return f, contenttype, nil
}

func (dp *DiskProvider) SaveFile(file io.Reader, filename string, path string) bool {
	rp := strings.Join([]string{dp.Location,path,filename}, "/")
	f, err := os.OpenFile(rp, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("error creating %v: %v\n", rp, err.Error())
		return false
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Printf("error writing %v: %v\n", rp, err.Error())
		return false
	}
	return true
}

func (dp *DiskProvider) ObjectInfo(path string) (bool, bool, string) {
	rp := strings.Join([]string{dp.Location,path}, "/")
	fileStat, err := os.Stat(rp)
	if err != nil {
		fmt.Printf("error gather stats for file %v: %v", rp, err.Error())
		return false, false, FileIsLocal
	}

	if fileStat.IsDir() {
		return true, true, FileIsLocal
	}
	return true, false, FileIsLocal
}

func (dp *DiskProvider) CreateDirectory(path string) bool {
	rp := strings.Join([]string{dp.Location,path}, "/")
	err := os.Mkdir(rp, 0755)
	if err != nil {
		fmt.Printf("Error creating directory %v: %v\n", rp, err.Error())
		return false
	}
	return true
}

func (dp *DiskProvider) Delete(path string) bool {
	rp := strings.Join([]string{dp.Location,path}, "/")
	err := os.RemoveAll(rp)
	if err != nil {
		fmt.Printf("Error removing %v: %v\n", path, err.Error())
		return false
	}
	return true
}