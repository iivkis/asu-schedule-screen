package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var groupID, groupToken string

func init() {
	var get = func(key string) string {
		s, ok := os.LookupEnv(key)
		if !ok {
			panic("undefined env: " + key)
		}
		return s
	}

	groupToken = get("ASU_SCHED_GP_TOKEN")
	groupID = get("ASU_SCHED_GP_ID")
}

func defaultRawQuery() string {
	q := url.Values{}
	q.Set("access_token", groupToken)
	q.Set("v", "5.100")
	q.Set("group_id", groupID)
	return q.Encode()
}

func getUploadUrl() (string, error) {
	var response struct {
		Response struct {
			UploadURL string `json:"upload_url"`
		} `json:"response"`
	}

	p, _ := url.Parse("https://api.vk.com/method/photos.getMessagesUploadServer")
	p.RawQuery = defaultRawQuery()

	res, err := http.Get(p.String())
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	b, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(b, &response); err != nil {
		return "", err
	}

	if response.Response.UploadURL == "" {
		return "", errors.New(string(b))
	}

	return response.Response.UploadURL, nil
}

func uploadAndSave(uploadURL string, photo []byte) (string, error) {
	var uploadResponse struct {
		Server int
		Hash   string
		Photo  string
	}

	//upload
	{
		var buf bytes.Buffer
		w := multipart.NewWriter(&buf)

		fw, err := w.CreateFormFile("photo", "schedule.png")
		if err != nil {
			return "", err
		}

		if _, err := io.Copy(fw, bytes.NewReader(photo)); err != nil {
			return "", err
		}
		w.Close()

		req, err := http.NewRequest("POST", uploadURL, &buf)
		if err != nil {
			return "", err
		}

		req.Header.Set("Content-Type", w.FormDataContentType())

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}
		defer res.Body.Close()

		b, _ := ioutil.ReadAll(res.Body)
		if err := json.Unmarshal(b, &uploadResponse); err != nil {
			return "", err
		}
	}

	var saveResponse struct {
		Response []struct {
			OwnerID int `json:"owner_id"`
			MediaID int `json:"id"`
		} `json:"response"`
	}
	//save
	{
		p, _ := url.Parse("https://api.vk.com/method/photos.saveMessagesPhoto")
		p.RawQuery = defaultRawQuery()

		q := p.Query()
		q.Set("server", strconv.Itoa(uploadResponse.Server))
		q.Set("hash", uploadResponse.Hash)
		q.Set("photo", uploadResponse.Photo)
		p.RawQuery = q.Encode()

		res, err := http.Get(p.String())
		if err != nil {
			return "", err
		}
		defer res.Body.Close()

		b, _ := ioutil.ReadAll(res.Body)
		if err := json.Unmarshal(b, &saveResponse); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("photo%d_%d", saveResponse.Response[0].OwnerID, saveResponse.Response[0].MediaID), nil
}
