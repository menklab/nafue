package utility

import (
	"nafue-api/models/display"
	"os"
	"regexp"
	"nafue/config"
	"nafue/models"
	"io/ioutil"
	"errors"
	"fmt"
	"strconv"
)

var fileIdRegex = regexp.MustCompile(`^.*file/(.*)$`)

func GetFile(url string) {

	// get api url from share link
	aUrl := appifyUrl(url)

	// download file header info
	var fileHeader = display.FileHeaderDisplay{}
	getFileHeader(aUrl, &fileHeader)

	// download file body
	secureFileBody := getFileBody(fileHeader.DownloadUrl)

	// loop until good pass or 3 attempts
	var fileBody *models.FileBody
	var err error
	var attemptCount = 0
	for attemptCount < 3 {
		// ask for password
		pass := promptPassword()

		// decrypt file with password
		fileBody, err = Decrypt(&fileHeader, pass, secureFileBody)

		// check for error and decide path
		if err == nil {
			break
		}
		fmt.Println("Failed to decrypt. Please try again.")
		attemptCount++
	}
	// write file to disk
	if err == nil {
		writeFileContentsToPath(fileBody)
		fmt.Println("File saved to: " + fileBody.Name)
	} else {
		fmt.Println("To many failed attempts. File was deleted.")
		os.Exit(1)
	}
}

func PutFile(file string) string {

	// ask for password
	pass := promptPassword()

	// get file contents
	fileBody := getFileContentsFromPath(file)

	// encrypt file with password
	secureData, fileHeader := Encrypt(fileBody, pass)

	// put file header info
	putFileHeader(config.API_FILE_URL, fileHeader)

	// post body data
	putFileBody(fileHeader.UploadUrl, secureData)

	// provide share link
	shareLink := config.SHARE_LINK + fileHeader.ShortUrl
	fmt.Println("Share Link: ", shareLink)
	return shareLink

}

func getFileContentsFromPath(path string) *models.FileBody {

	// verify file is under 50mb
	fileInfo, err := os.Stat(path)
	checkError(err)
	fileSize := fileInfo.Size()
	if fileSize > (config.FILE_SIZE_LIMIT * 1024 * 1024) {
		checkError(errors.New("File is larger than " + strconv.FormatInt(config.FILE_SIZE_LIMIT, 10) + "mb."))
	}

	// get file type and name
	fileName := fileInfo.Name()

	// read file
	fileBytes, err := ioutil.ReadFile(path)
	checkError(err)

	// create file data package
	fbp := models.FileBody{
		Name: fileName,
		Part: 0,
		Content: fileBytes,
	}
	return &fbp
}

func writeFileContentsToPath(fileBody *models.FileBody) {

	err := ioutil.WriteFile(fileBody.Name, fileBody.Content, 0644)
	checkError(err)

}

