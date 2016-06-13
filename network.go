package nafue

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/menkveldj/nafue-api/models/display"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ErrorDisplay struct {
	Message string `json:"message,omitempty"`
}

func getFileHeader(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		errorDisplay := ErrorDisplay{}
		err := json.NewDecoder(r.Body).Decode(&errorDisplay)
		if err != nil {
			return err
		}
		return errors.New(strconv.Itoa(r.StatusCode) + ", " + errorDisplay.Message)
	}

	return json.NewDecoder(r.Body).Decode(target)
}

func getFileBody(url string) (*[]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read body
	rBody, err := ioutil.ReadAll(resp.Body)
	if (err != nil) {
		return nil, err
	}
	return &rBody, nil
}

func putFileHeader(url string, fileHeader *display.FileHeaderDisplay) error {
	// create json body
	body, err := json.Marshal(&fileHeader)

	// create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if (err != nil) {
		return nil
	}
	req.Header.Set("Content-Type", "application/json")

	// make client and do request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	rBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(rBody, &fileHeader)
	if (err != nil) {
		return err
	}

	return nil
}

func putFileBody(url string, body *[]byte) error {
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(*body))
	if (err != nil) {
		return err
	}

	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")

	// make client and do request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func appifyUrl(url string) string {
	fileId := fileIdRegex.FindStringSubmatch(url)[1]
	// use fileId to get file from api
	appifiedUrl := C.API_URL + "/files/" + fileId
	return appifiedUrl
}
