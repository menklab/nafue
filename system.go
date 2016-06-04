package nafue

import "os"
import (
	"fmt"
	"github.com/menkveldj/nafue/config"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"syscall"
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

func promptPassword() string {
	// ask for password
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	checkError(err)

	password := string(bytePassword)
	fmt.Println()
	return password
}
