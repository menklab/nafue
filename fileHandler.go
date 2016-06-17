package nafue

import (
	"errors"
	"io"
	"regexp"
	"strconv"
	"os"
	"crypto/sha256"
	"github.com/menkveldj/nafue-api/models/display"
)

var fileIdRegex = regexp.MustCompile(`^.*file/(.*)$`)


func GetFile(url string, secureData io.ReadWriteSeeker) (*display.FileHeaderDisplay, error) {

	// get api url from share link
	aUrl := appifyUrl(url)

	// download file header info
	fileHeader, err := getFileHeader(aUrl)
	if err != nil {
		return nil, err
	}

	// download file body
	err = getFileBody(secureData, fileHeader)
	if err != nil {
		return nil, err
	}

	return fileHeader, nil
}

func UnsealFile (secureData io.ReadWriteSeeker, pass string, fileHeader *display.FileHeaderDisplay, fileInfo os.FileInfo) error {

	// decrypt to file
	err := Decrypt(secureData, pass, fileHeader)
	if err != nil {
		return err
	}

	return nil
}

func SealShareFile(data io.ReaderAt, secureData io.ReadWriteSeeker, fileInfo os.FileInfo, name, pass string) (string, error) {

	// check file is under 50mb
	if fileInfo.Size() > (C.FILE_SIZE_LIMIT * 1024 * 1024) {
		err := errors.New("File is larger than " + strconv.FormatInt(C.FILE_SIZE_LIMIT, 10) + "mb.")
		return "", err
	}

	// encrypt to temp file
	fileHeader, err := Encrypt(data, secureData, fileInfo.Name(), pass)
	if err != nil {
		return "", err
	}

	// set reader to start of file
	_, err = secureData.Seek(0,0)
	if err != nil {
		return "", err
	}

	 //create checksum
	checksum := sha256.New()
	fileSize, err := io.Copy(checksum, secureData)
	if err != nil {
		return "", nil
	}
	fileHeader.MD5Checksum = checksum.Sum(nil)
	fileHeader.FileSize = fileSize

	err = putFileHeader(C.API_FILE_URL, fileHeader)
	if err != nil {
		return "", errors.New("PutFileHeader: " + err.Error())
	}

	// post body data
	err = putFileBody(fileHeader, secureData)
	if err != nil {
		return "", errors.New("PutFileBody: " + err.Error())
	}

	// provide share link
	shareLink := C.SHARE_LINK + fileHeader.ShortUrl
	return shareLink, nil
}

//
//func getFileContentsFromReader(reader io.Reader, size int64, name string) (*models.FileBody, error) {
//
//
//	fileBytes, err := ioutil.ReadAll(reader)
//	if err != nil {
//		return nil , err
//	}
//
//	fbp := models.FileBody{
//		Name:    name,
//		Part:    0,
//		Content: fileBytes,
//	}
//
//	return &fbp, nil
//}
//
//func writeFileContentsToPath(fileBody *models.FileBody) error {
//
//	err := ioutil.WriteFile(fileBody.Name, fileBody.Content, 0644)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
