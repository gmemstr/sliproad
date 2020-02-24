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
}

type BackblazeAuthPayload struct {
	AccountId string `json:"accountId"`
	AuthToken string `json:"authorizationToken"`
	ApiUrl string `json:"apiUrl"`
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
	bp.Name = "Backblaze|" + data.AccountId
	
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

func (bp *BackblazeProvider) ViewFile(path string) string {
	return ""
}

func (bp *BackblazeProvider) SaveFile(contents []byte, path string) bool {
	return true
}