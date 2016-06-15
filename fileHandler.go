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


func UnSealFile(url string) (*[]byte, *display.FileHeaderDisplay, error) {
//
//	// get api url from share link
//	aUrl := appifyUrl(url)
//
//	// dowload file header info
//	var fileHeader = display.FileHeaderDisplay{}
//	getFileHeader(aUrl, &fileHeader)
//
//	// dowload file body
//	fileBody, err := getFileBody(fileHeader.DownloadUrl)
//	if err != nil {
//		return nil, nil, err
//	}
//	return fileBody, &fileHeader, nil
	return nil, nil, nil
}
//
//func TryDecrypt(body *[]byte, header *display.FileHeaderDisplay, pass string) (io.Reader, string, error) {
//	fileBody, err := Decrypt(header, pass, body)
//	if err != nil {
//		return bytes.NewBufferString(""), "", err
//	}
//	return bytes.NewBuffer(fileBody.Content), fileBody.Name, nil
//}

func SealFile(data io.ReaderAt, secureData io.ReadWriteSeeker, fileInfo os.FileInfo, name, pass string) (string, error) {

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

	// create checksum
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
