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
	aUrl := appifyUrl(url)

	// download file header info
	var fileHeader = display.FileDisplay{}
	getFileHeader(aUrl, &fileHeader)

	// download file body
	secureData := getFileBodyByUrl(fileHeader.DownloadUrl)
	log.Println("File Body: ", *secureData)

	// ask for password
	pass := promptPassword()

	// decrypt file with password
	fileBody := Decrypt(&fileHeader, pass, secureData)
	log.Println("Decrypted Body: ", fileBody)
}

func PutFile(file string) string{

	// ask for password
	pass := promptPassword()

	// encrypt file with password
	secureData, fileHeader := Encrypt(file, pass)

	// put file header info
	putFileHeader(config.API_FILE_URL, fileHeader)

	// post body data
	putFileBody(fileHeader.UploadUrl, secureData)

	// provide share link
	shareLink := config.SHARE_LINK + fileHeader.ShortUrl
	log.Println("Share Link: ", shareLink)
	return shareLink

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

