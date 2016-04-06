package config

import (
	"path/filepath"
	"os"
	"crypto/sha1"
)


var (
	API_URL string = API_PROTOCOL + "://" + API_HOST + ":" + API_PORT + "/" + API_BASE
	API_FILE_URL string = API_URL + "/files"
	HASH_TYPE = sha1.New

)
const (
	API_PROTOCOL string = "http"
	API_HOST string = "localhost"
	API_PORT string = "9090"
	API_BASE string = "api"
	TEMP_DIR string = "nafue"
	SHARE_LINK string = "http://localhost:8080/file/"
	ITERATIONS int = 1000
	KEY_LENGTH int = 32
	SALT_LENGTH int = 32
	FILE_SIZE_LIMIT int64 = 50 // 50 mb
	NAFUE_TEMP_FILE string = ".tmp.nafue"

)

func GetTempDir() string {
	return filepath.Join(os.Getenv("HOME"), TEMP_DIR)
}
