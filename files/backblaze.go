package files

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	Size int `json:"contentLength"`
	Type string `json:"contentType"`
	FileName string `json:"fileName"`
	Timestamp int `json:"uploadTimestamp"`
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

// Call Backblaze API endpoint to authorize and gather facts.
func (bp *BackblazeProvider) Authorize(appKeyId string, appKey string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET",
		"https://api.backblazeb2.com/b2api/v2/b2_authorize_account",
		nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(appKeyId, appKey)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data BackblazeAuthPayload

	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	bp.Authentication = data.AuthToken
	bp.Location = data.ApiUrl
	bp.Name = data.AccountId
	bp.DownloadLocation = data.DownloadUrl
	
	return nil
}

func (bp *BackblazeProvider) GetDirectory(path string) Directory {
	client := &http.Client{}
	
	requestBody := fmt.Sprintf(`{"bucketId": "%s"}`, bp.Bucket)

	req, err := http.NewRequest("POST",
		bp.Location + "/b2api/v2/b2_list_file_names",
		bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return Directory{}
	}
	req.Header.Add("Authorization", bp.Authentication)
	resp, err := client.Do(req)
	if err != nil {
		return Directory{}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
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

func (bp *BackblazeProvider) ViewFile(path string) []byte {
	client := &http.Client{}
	// Get bucket name >:(
	bucketIdPayload := fmt.Sprintf(`{"accountId": "%s", "bucketId": "%s"}`, bp.Name, bp.Bucket)

	req, err := http.NewRequest("POST", bp.Location + "/b2api/v2/b2_list_buckets",
		bytes.NewBuffer([]byte(bucketIdPayload)))
	req.Header.Add("Authorization", bp.Authentication)

	res, err := client.Do(req)
	bucketData, err := ioutil.ReadAll(res.Body)

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
		return nil
	}
	req.Header.Add("Authorization", bp.Authentication)
	file, err := client.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	fileBytes, err := ioutil.ReadAll(file.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return fileBytes
}

func (bp *BackblazeProvider) SaveFile(contents []byte, path string) bool {
	return true
}

func (bp *BackblazeProvider) DetermineType(path string) string {
	return "file"
}