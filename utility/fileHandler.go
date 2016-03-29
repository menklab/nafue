package utility

import (
	"nafue-api/models/display"
	"log"
	"os"
	"regexp"
	"nafue/config"
	"path/filepath"
)

var fileIdRegex = regexp.MustCompile(`^.*file/(.*)$`)

func GetFile(url string) {

	// get fileId
	aUrl := appifyUrl(url)

	// download file header info
	var fileHeader = display.FileDisplay{}
	getFileHeader(aUrl, &fileHeader)

	// download file body
	getFileBody(&fileHeader)

}

func appifyUrl(url string) string {
	fileId := fileIdRegex.FindStringSubmatch(url)[1]
	// use fileId to get file from api
	appifiedUrl := config.API_URL + "/files/" + fileId
	log.Println("appifiedUrl: ", appifiedUrl)
	return appifiedUrl
}

func getFileHeader(url string, fileHeader *display.FileDisplay){
	err := getJson(url, &fileHeader)
	if err != nil {
		log.Println("Error retrieving file header: ", err.Error())
		os.Exit(1)
	}
	log.Println("file info: ", fileHeader.ToString())
}

func getFileBody(fileHeader *display.FileDisplay) {
	err := getFileByUrl(filepath.Join(config.GetTempDir(), fileHeader.ShortUrl), fileHeader.DownloadUrl)
	if err != nil {
		log.Println("Error retrieving file body: ", err.Error())
		os.Exit(1)
	}
	return
}