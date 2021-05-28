package files

import (
	"fmt"
	"io"
	// I _really_ don't want to deal with AWS API stuff by hand.
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var svc s3.S3

type S3Provider struct {
	FileProvider
	Region string
	Bucket string
}


// Setup runs when the application starts up, and allows for things like authentication.
func (s *S3Provider) Setup(args map[string]string) bool {
	sess, err := session.NewSession(&aws.Config{
    	Region: aws.String(s.Region)},
	)
	if err != nil {
		return false
	}
	svc = *s3.New(sess)
	return true
}

// GetDirectory fetches a directory's contents.
func (s *S3Provider) GetDirectory(path string) Directory {
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(s.Bucket)})
	if err != nil {
		fmt.Println(err)
	    return Directory{}
	}

	dir := Directory{}
	for _, item := range resp.Contents {
		file := FileInfo{
			IsDirectory: false,
			Name: *item.Key,
		}
		dir.Files = append(dir.Files, file)
	}

	return dir
}

// RemoteFile will bypass http.ServeContent() and instead write directly to the response.
func (s *S3Provider) SendFile(path string, writer io.Writer) (stream io.Reader, contenttype string, err error) {
	return
}

// SaveFile will save a file with the contents of the io.Reader at the path specified.
func (s *S3Provider) SaveFile(file io.Reader, filename string, path string) bool {
	return false
}

// ObjectInfo will return the info for an object given a path to if the file exists and location.
// Should return whether the path exists, if the path is a directory, and if it lives on disk.
// (see constants defined: `FILE_IS_REMOTE` and `FILE_IS_LOCAL`)
func (s *S3Provider) ObjectInfo(path string) (bool, bool, string) {
	return true, true, ""
}

// CreateDirectory will create a directory on services that support it.
func (s *S3Provider) CreateDirectory(path string) bool {
	return false
}

// Delete simply deletes a file. This is expected to be a destructive action by default.
func (s *S3Provider) Delete(path string) bool {
	return false
}
