package utility

import (
	"nafue-api/models/display"
	"log"
	"os"
	"regexp"
	"nafue/config"
)

var fileIdRegex = regexp.MustCompile(`^.*file/(.*)$`)

func GetFile(url string) {

	// get fileId
	appifyUrl(url)
	// first download file
	//getFileHeader(url)
}

func appifyUrl(url string) {
	// http://localhost:8080/file/a5627316-2805-4cec-7700-5dceb2df0911
	fileId := fileIdRegex.FindStringSubmatch(url)[1]
	// use fileId to get file from api
	appifiedUrl := config.API_URL + "/files/" + fileId
	log.Println("appifiedUrl: " + appifiedUrl)
	getFileHeader(appifiedUrl)
}

func getFileHeader(url string) {
	fileHeader := new(display.FileDisplay)
	err := getJson(url, &fileHeader)
	if err != nil {
		log.Println("Error retrieving file: ", err.Error())
		os.Exit(1)
	}
	log.Println("file info: ", fileHeader.ToString())
}
