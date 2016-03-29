package utility

import (
	"encoding/json"
	"net/http"
	"log"
	"errors"
	"strconv"
	"os"
	"io"
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

func getFileByUrl(tempFile string, url string) error {

	// Create the file
	out, err := os.Create(tempFile)
	if err != nil {
		log.Println("Error creating file: ", tempFile)
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Erorr downloading file data: " + url)
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Println("Erorr writing data to file: ", tempFile)

		return err
	}

	return nil
}
