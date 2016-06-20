package nafue

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/menkveldj/nafue-api/models/display"
	"io/ioutil"
	"net/http"
	"io"
	"os"
)

func getFileHeader(url string) (*display.FileHeaderDisplay, error) {
	r, err := http.Get(url)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}

	if r.StatusCode != 200 {
		return nil, errors.New("GetFileHeader: Services responded with " + r.Status)
	}

	fileHeader := display.FileHeaderDisplay{}
	err = json.NewDecoder(r.Body).Decode(&fileHeader)
	if err != nil {
		return nil, err
	}

	return &fileHeader, nil
}

func getFileBody(secureData io.ReadWriteSeeker, fileHeader *display.FileHeaderDisplay) error {

	r, err := http.Get(fileHeader.DownloadUrl)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return errors.New("GetFileBody: Services responded with " + r.Status)
	}

	// read body
	_, err = io.Copy(secureData, r.Body)
	if err != nil {
		return err
	}

	return nil
}

func putFileHeader(url string, fileHeader *display.FileHeaderDisplay) error {
	// create json body
	body, err := json.Marshal(&fileHeader)

	// create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if (err != nil) {
		return err
	}

	// make client and do request
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("Services responded with " + resp.Status)
	}

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

func putFileBody(fileHeader *display.FileHeaderDisplay, secureFile *os.File) error {

	// make sure we read file form start
	_,  err := secureFile.Seek(0,0)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fileHeader.UploadUrl, secureFile)
	if err != nil {
		return err
	}

	// get file size
	fStat, err := secureFile.Stat()
	if err != nil {
		return err
	}
	req.ContentLength = fStat.Size()
	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")

	// make client and do request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return err
	//}

	if resp.StatusCode != http.StatusOK {
		return errors.New("Services responded with " + resp.Status)
	}

	return nil
}

func appifyUrl(url string) string {
	fileId := fileIdRegex.FindStringSubmatch(url)[1]
	// use fileId to get file from api
	appifiedUrl := C.API_URL + "/files/" + fileId
	return appifiedUrl
}
