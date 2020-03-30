package files

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
)

type BackblazeProvider struct {
	FileProvider
	Bucket string
	DownloadLocation string
}

type BackblazeAuthPayload struct {
	AccountId string `json:"accountId"`
	AuthToken string `json:"authorizationToken"`
	ApiUrl string `json:"apiUrl"`
	DownloadUrl string `json:"downloadUrl"`
}

type BackblazeFile struct {
	Action string `json:"action"`
	Size int64 `json:"contentLength"`
	Type string `json:"contentType"`
	FileName string `json:"fileName"`
	Timestamp int64 `json:"uploadTimestamp"`
}

type BackblazeFilePayload struct {
	Files []BackblazeFile `json:"files"`
}

type BackblazeBucketInfo struct {
	BucketId string `json:"bucketId"`
	BucketName string `json:"bucketName"`
}

type BackblazeBucketInfoPayload struct {
	Buckets []BackblazeBucketInfo `json:"buckets"`
}

type BackblazeUploadInfo struct {
	UploadUrl string `json:"uploadUrl"`
	AuthToken string `json:"authorizationToken"`
}

// Call Backblaze API endpoint to authorize and gather facts.
func (bp *BackblazeProvider) Setup(args map[string]string) bool {
	applicationKeyId := args["applicationKeyId"]
	applicationKey := args["applicationKey"]

	client := &http.Client{}
	req, err := http.NewRequest("GET",
		"https://api.backblazeb2.com/b2api/v2/b2_authorize_account",
		nil)
	if err != nil {
		return false
	}
	req.SetBasicAuth(applicationKeyId, applicationKey)
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return false
	}

	var data BackblazeAuthPayload

	err = json.Unmarshal(body, &data)
	if err != nil {
		return false
	}

	bp.Authentication = data.AuthToken
	bp.Location = data.ApiUrl
	bp.Name = data.AccountId
	bp.DownloadLocation = data.DownloadUrl

	return true
}

func (bp *BackblazeProvider) GetDirectory(path string) Directory {
	client := &http.Client{}
	
	requestBody := fmt.Sprintf(`{"bucketId": "%s"}`, bp.Bucket)

	req, err := http.NewRequest("POST",
		bp.Location + "/b2api/v2/b2_list_file_names",
		bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		fmt.Println(err.Error())
		return Directory{}
	}
	req.Header.Add("Authorization", bp.Authentication)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return Directory{}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return Directory{}
	}

	var data BackblazeFilePayload
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err.Error())
		return Directory{}
	}
	finalDir := Directory{
		Path: bp.Bucket,
	}
	for _, v := range data.Files {
		file := FileInfo{
			IsDirectory: v.Action == "folder",
			Name:        v.FileName,
		}
		if v.Action != "folder" {
			split := strings.Split(v.FileName, ".")
			file.Extension = split[len(split) - 1]
		}
		finalDir.Files = append(finalDir.Files, file)
	}

	return finalDir
}

func (bp *BackblazeProvider) ViewFile(path string, w io.Writer) {
	client := &http.Client{}
	// Get bucket name >:(
	bucketIdPayload := fmt.Sprintf(`{"accountId": "%s", "bucketId": "%s"}`, bp.Name, bp.Bucket)

	req, err := http.NewRequest("POST", bp.Location + "/b2api/v2/b2_list_buckets",
		bytes.NewBuffer([]byte(bucketIdPayload)))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	req.Header.Add("Authorization", bp.Authentication)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	bucketData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var data BackblazeBucketInfoPayload
	json.Unmarshal(bucketData, &data)
	ourBucket := data.Buckets[0].BucketName
	// Get file and write over to client.
	url := strings.Join([]string{bp.DownloadLocation,"file",ourBucket,path}, "/")
	req, err = http.NewRequest("GET",
		url,
		bytes.NewBuffer([]byte("")))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	req.Header.Add("Authorization", bp.Authentication)
	file, err := client.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = io.Copy(w, file.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (bp *BackblazeProvider) SaveFile(file multipart.File, handler *multipart.FileHeader, path string) bool {
	client := &http.Client{}
	bucketIdPayload := fmt.Sprintf(`{"bucketId": "%s"}`, bp.Bucket)

	req, err := http.NewRequest("POST", bp.Location + "/b2api/v2/b2_get_upload_url",
		bytes.NewBuffer([]byte(bucketIdPayload)))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	req.Header.Add("Authorization", bp.Authentication)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	bucketData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	var data BackblazeUploadInfo
	json.Unmarshal(bucketData, &data)

	req, err = http.NewRequest("POST",
		data.UploadUrl,
		file,
	)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	// Read the content
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(req.Body)
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	// Calculate SHA1 and add required headers.
	fileSha := sha1.New()
	fileSha.Write(bodyBytes)

	req.Header.Add("Authorization", data.AuthToken)
	req.Header.Add("X-Bz-File-Name", handler.Filename)
	req.Header.Add("Content-Type", "b2/x-auto")
	req.Header.Add("X-Bz-Content-Sha1", fmt.Sprintf("%x", fileSha.Sum(nil)))
	req.ContentLength = handler.Size

	// Upload in background.
	go func() {
		res, err = client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	return true
}

func (bp *BackblazeProvider) DetermineType(path string) string {
	return "file"
}

func (bp *BackblazeProvider) GetProviderConfig() FileProvider {
	return ProviderConfig[bp.Name]
}