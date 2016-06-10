package nafue

import (
	"os"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"syscall"
	"github.com/menkveldj/nafue/config"
)

var C config.Config

func Init(c config.Config) {

	// init config
	C = c

	// check temp
	setupTemp()
}

func setupTemp() {

	err := os.MkdirAll(C.TEMP_DIR, 0700)
	if err != nil {
		log.Println("error creating temp directory: ", err.Error())
	}
}

func promptPassword() string {
	// ask for password
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	checkError(err)

	password := string(bytePassword)
	fmt.Println()
	return password
}

func checkError(e error) {
	if e != nil {
		fmt.Println("Error: ", e.Error())
		os.Exit(1)
	}
}

