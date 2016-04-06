package utility

import (
	"encoding/json"
	"net/http"
	"log"
	"errors"
	"strconv"
	"os"
	"io/ioutil"
	"bytes"
	"nafue/config"
	"nafue-api/models/display"
)

type ErrorDisplay struct {
	Message string `json:"message,omitempty"`
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		errorDisplay := ErrorDisplay{}
		err := json.NewDecoder(r.Body).Decode(&errorDisplay)
		if (err != nil) {
			log.Println("Error getting service, error message: ", err.Error())
			return err
		}
		return errors.New(strconv.Itoa(r.StatusCode) + ", " + errorDisplay.Message)
	}

	return json.NewDecoder(r.Body).Decode(target)
}

func getFileBodyByUrl(url string) *string {

	resp, err := http.Get(url)
	if err != nil {
		checkError(err)
	}
	defer resp.Body.Close()

	// read body
	rBody, err := ioutil.ReadAll(resp.Body)
	checkError(err)
	rBodyAsString := string(rBody)
	return &rBodyAsString
}

func putFileHeader(url string, fileHeader *display.FileDisplay) {
	log.Println("\n\n: ########## Put File Header ##########")
	log.Println("POST: " + url)
	// create json body
	body, err := json.Marshal(&fileHeader)

	log.Println("body: ", string(body))
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

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	rBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error in service response: ", err.Error())
		os.Exit(1)
	}

	log.Println("response Body:", string(rBody))

	err = json.Unmarshal(rBody, &fileHeader)
	checkError(err)
}

func putFileBody(url string, body *string) {
	log.Println("\n\n: ########## Put File Body ##########")
	log.Println("PUT: ", url)
	log.Println("file body: ", *body)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(*body)))
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

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	rBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error in service response: ", err.Error())
		os.Exit(1)
	}

	log.Println("response Body:", string(rBody))

}

func appifyUrl(url string) string {
	fileId := fileIdRegex.FindStringSubmatch(url)[1]
	// use fileId to get file from api
	appifiedUrl := config.API_URL + "/files/" + fileId
	return appifiedUrl
}
