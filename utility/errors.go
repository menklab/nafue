package utility

import (
	"os"
	"log"
)

func checkError(e error) {
	if e != nil {
		log.Println("Error: ", e.Error())
		os.Exit(1)
	}
}
