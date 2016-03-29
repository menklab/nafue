package utility

import "os"
import (
	"nafue/config"
	"log"
)

func Init() {

	// check temp
	setupTemp()
}

func setupTemp() {

	err := os.MkdirAll(config.GetTempDir(), 0700)
	if err != nil {
		log.Println("error creating temp directory: ", err.Error())
	}
}
