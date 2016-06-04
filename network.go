package nafue

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/menkveldj/nafue-api/models/display"
	"github.com/menkveldj/nafue/config"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
			log.Println("Error getting service, error message: ", err.Error())
			return err
		}
		return errors.New(strconv.Itoa(r.StatusCode) + ", " + errorDisplay.Message)
	}

	return json.NewDecoder(r.Body).Decode(target)
}

func getFileBody(url string) *[]byte {

	resp, err := http.Get(url)
	if err != nil {
		checkError(err)
	}
	defer resp.Body.Close()

	// read body
	rBody, err := ioutil.ReadAll(resp.Body)
	checkError(err)
	return &rBody
}

func putFileHeader(url string, fileHeader *display.FileHeaderDisplay) {
	// create json body
	body, err := json.Marshal(&fileHeader)

	// create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	checkError(err)
	req.Header.Set("Content-Type", "application/json")

	// make client and do request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error posting fileheader data: ", err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	rBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error in service response: ", err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(rBody, &fileHeader)
	checkError(err)
}

func putFileBody(url string, body *[]byte) {
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(*body))
	checkError(err)
	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")

	// make client and do request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error posting fileheader data: ", err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error in service response: ", err.Error())
		os.Exit(1)
	}

}

func appifyUrl(url string) string {
	fileId := fileIdRegex.FindStringSubmatch(url)[1]
	// use fileId to get file from api
	appifiedUrl := config.API_URL + "/files/" + fileId
	return appifiedUrl
}
