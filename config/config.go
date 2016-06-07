package config

import (
	"path/filepath"
	"os"
	"crypto/sha1"
	"hash"
)

type Config struct {
	API_PROTOCOL    string
	API_HOST        string
	API_PORT        string
	API_BASE        string
	TEMP_DIR        string
	SHARE_LINK      string
	ITERATIONS      int
	KEY_LENGTH      int
	SALT_LENGTH     int
	FILE_SIZE_LIMIT int
	NAFUE_TEMP_FILE string
	API_URL         string
	API_FILE_URL    string
	HASH_TYPE       hash.Hash
}

var (
	Current Config
)

// get production configuration for lib
func Production() Config {
	c := Config{
		API_PROTOCOL: "https",
		API_HOST: "api.nafue.com",
		API_PORT: "80",
		API_BASE: "api",
		TEMP_DIR: filepath.Join(os.Getenv("HOME"), "nafue"),
		SHARE_LINK: "https://www.nafue.com/file/",
		ITERATIONS: 1000,
		KEY_LENGTH: 32,
		SALT_LENGTH: 32,
		FILE_SIZE_LIMIT: 50, // 50 mb
		NAFUE_TEMP_FILE: ".tmp.nafue",
		HASH_TYPE: sha1.New(),
	}
	c.API_URL = Current.API_PROTOCOL + "://" + Current.API_HOST + ":" + Current.API_PORT + "/" + Current.API_BASE;
	c.API_FILE_URL = Current.API_URL + "/files";
	return c
}

// set a cost configuration for the lib
func Set(c Config) {
	Current = c
}

// get current configuration utilized by lib
func Get() Config {
	return Current
}
