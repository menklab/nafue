package utility

import (
	"encoding/json"
	"net/http"
	"log"
	"errors"
	"strconv"
)

type ErrorDisplay struct {
	Message    string `json:"message,omitempty"`
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
			log.Println("Error getting service error message: ", err.Error())
			return err
		}
		return errors.New(strconv.Itoa(r.StatusCode) + ", " + errorDisplay.Message)
	}

	return json.NewDecoder(r.Body).Decode(target)
}
