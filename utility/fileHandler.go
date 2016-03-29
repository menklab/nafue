package utility

import (
	"nafue-api/models/display"
	"log"
	"os"
	"regexp"
	"nafue/config"
	"path/filepath"
	//b64 "encoding/base64"
	"encoding/json"
"net/http"
	"bytes"
	"io/ioutil"
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

func PutFile(file string, password string) {

	// encrypt file with password
	Encrypt(password)
	// mode=ccm, cipher=aes, tag=128, key=256, iterations=1000
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
	// decode salt
	//saltByte, err := b64.StdEncoding.DecodeString(fileHeader.Salt)
	//if err != nil {
	//	log.Println("Error decoding salt: ", err.Error())
	//	os.Exit(1)
	//}

	// decode iv
	//iv, err := b64.StdEncoding.DecodeString(fileHeader.IV)
	//if err != nil {
	//	log.Println("Error decoding salt: ", err.Error())
	//	os.Exit(1)
	//}
}

func getFileBody(fileHeader *display.FileDisplay) {
	err := getFileBodyByUrl(filepath.Join(config.GetTempDir(), fileHeader.ShortUrl), fileHeader.DownloadUrl)
	if err != nil {
		log.Println("Error retrieving file body: ", err.Error())
		os.Exit(1)
	}
	return
}

func putFileHeader(url string, fileHeader *display.FileDisplay) {
	// create json body
	body, err := json.Marshal(&fileHeader)

	// create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
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

	log.Println("response Body:", rBody)
}