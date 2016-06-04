package nafue

import (
	"fmt"
	"os"
)

func checkError(e error) {
	if e != nil {
		fmt.Println("Error: ", e.Error())
		os.Exit(1)
	}
}
