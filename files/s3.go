package files

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"

	// I _really_ don't want to deal with AWS API stuff by hand.
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var svc s3.S3
var sess *session.Session

type S3Provider struct {
	FileProvider
	Region    string
	Bucket    string
	Endpoint  string
	KeyID     string
	KeySecret string
}

// Setup runs when the application starts up, and allows for things like authentication.
func (s *S3Provider) Setup(args map[string]string) bool {
	config := &aws.Config{Region: aws.String(s.Region)}
	if s.KeyID != "" && s.KeySecret != "" {
		config = &aws.Config{
			Region:      aws.String(s.Region),
			Credentials: credentials.NewStaticCredentials(s.KeyID, s.KeySecret, ""),
		}
	}
	if s.Endpoint != "" {
		config.Endpoint = &s.Endpoint
	}
	sess, err := session.NewSession(config)
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
		ik := *item.Key
		// Why is this here? AWS returns a complete list of files, including
		// files within subdirectories (prefixed with the dir name). So we can
		// ignore directories altogether -- I would prefer to display them but
		// not sure what the best method of distinguishing them in ObjectInfo()
		// would be.
		if ik[len(ik)-1:] == "/" {
			continue
		}
		file := FileInfo{
			IsDirectory: false,
			Name:        *item.Key,
		}
		dir.Files = append(dir.Files, file)
	}

	return dir
}

func (s *S3Provider) SendFile(path string) (stream io.Reader, contenttype string, err error) {
	req, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &path,
	})
	if err != nil {
		return stream, contenttype, err
	}

	contenttype = mime.TypeByExtension(filepath.Ext(path))
	if contenttype == "" {
		var buf [512]byte
		n, _ := io.ReadFull(req.Body, buf[:])
		contenttype = http.DetectContentType(buf[:n])
	}

	return req.Body, contenttype, err
}

func (s *S3Provider) SaveFile(file io.Reader, filename string, path string) bool {
	return false
}

func (s *S3Provider) ObjectInfo(path string) (bool, bool, string) {
	if path == "" {
		return true, true, ""
	}

	_, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &path,
	})
	if err != nil {
		fmt.Println(err)
		return false, false, ""
	}

	return true, false, ""
}

// CreateDirectory will create a directory on services that support it.
func (s *S3Provider) CreateDirectory(path string) bool {
	return false
}

// Delete simply deletes a file. This is expected to be a destructive action by default.
func (s *S3Provider) Delete(path string) bool {
	return false
}
