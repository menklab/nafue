package nafue

import (
	"os"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"syscall"
	"nafue/config"
)

func Init(c config.Config) {

	// init config
	config.Set(c)

	// check temp
	setupTemp()
}

func setupTemp() {

	err := os.MkdirAll(config.Current.TEMP_DIR, 0700)
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

