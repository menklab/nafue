package config

import (
	"path/filepath"
	"os"
)


var (
	API_URL string = API_PROTOCOL + "://" + API_HOST + ":" + API_PORT + "/" + API_BASE
)
const (
	API_PROTOCOL string = "http"
	API_HOST string = "localhost"
	API_PORT string = "9090"
	API_BASE string = "api"
	TEMP_DIR string = "nafue"
)

func GetTempDir() string {
	return filepath.Join(os.Getenv("HOME"), TEMP_DIR)
}
