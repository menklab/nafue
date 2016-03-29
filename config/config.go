package config

import (
	"path/filepath"
	"os"
	"crypto/sha1"
	//"hash"
)


var (
	API_URL string = API_PROTOCOL + "://" + API_HOST + ":" + API_PORT + "/" + API_BASE
	HASH_TYPE = sha1.New

)
const (
	API_PROTOCOL string = "http"
	API_HOST string = "localhost"
	API_PORT string = "9090"
	API_BASE string = "api"
	TEMP_DIR string = "nafue"
	ITERATIONS int = 1000
	KEY_LENGTH int = 32
	SALT_LENGTH int = 32
)

func GetTempDir() string {
	return filepath.Join(os.Getenv("HOME"), TEMP_DIR)
}
