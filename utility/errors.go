package utility

import (
	"os"
	"fmt"
)

func checkError(e error) {
	if e != nil {
		fmt.Println("Error: ", e.Error())
		os.Exit(1)
	}
}
