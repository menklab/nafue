package config

import (
	"path/filepath"
	"os"
	"crypto/sha1"
	"hash"
)

type Config struct {
	API_PROTOCOL        string
	API_HOST            string
	API_PORT            string
	API_BASE            string
	TEMP_DIR            string
	SHARE_LINK          string
	ITERATIONS          int
	KEY_LENGTH          int
	SALT_LENGTH         int64
	FILE_SIZE_LIMIT     int64
	NAFUE_TEMP_FILE     string
	API_URL             string
	API_FILE_URL        string
	HASH_TYPE           func() hash.Hash
	BUFFER_SIZE         int64
	MAX_FILENAME_LENGTH int64
}

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
		HASH_TYPE: sha1.New,
		BUFFER_SIZE: 3200,
		MAX_FILENAME_LENGTH: 255,
	}
	c.API_URL = c.API_PROTOCOL + "://" + c.API_HOST + ":" + c.API_PORT + "/" + c.API_BASE;
	c.API_FILE_URL = c.API_URL + "/files";
	return c
}

func Development() Config {
	c := Config{
		API_PROTOCOL: "http",
		API_HOST: "dev-api.nafue.com",
		API_PORT: "80",
		API_BASE: "api",
		TEMP_DIR: filepath.Join(os.Getenv("HOME"), "nafue"),
		SHARE_LINK: "http://dev.nafue.com/file/",
		ITERATIONS: 1000,
		KEY_LENGTH: 32,
		SALT_LENGTH: 32,
		FILE_SIZE_LIMIT: 50, // 50 mb
		NAFUE_TEMP_FILE: ".tmp.nafue",
		HASH_TYPE: sha1.New,
		BUFFER_SIZE: 3200,
		MAX_FILENAME_LENGTH: 255,
	}
	c.API_URL = c.API_PROTOCOL + "://" + c.API_HOST + ":" + c.API_PORT + "/" + c.API_BASE;
	c.API_FILE_URL = c.API_URL + "/files";
	return c
}

func Local() Config {
	c := Config{
		API_PROTOCOL: "http",
		API_HOST: "localhost",
		API_PORT: "9090",
		API_BASE: "api",
		TEMP_DIR: filepath.Join(os.Getenv("HOME"), "nafue"),
		SHARE_LINK: "http://localhost/file/",
		ITERATIONS: 1000,
		KEY_LENGTH: 32,
		SALT_LENGTH: 32,
		FILE_SIZE_LIMIT: 50, // 50 mb
		NAFUE_TEMP_FILE: ".tmp.nafue",
		HASH_TYPE: sha1.New,
		BUFFER_SIZE: 3200,
		MAX_FILENAME_LENGTH: 255,
	}
	c.API_URL = c.API_PROTOCOL + "://" + c.API_HOST + ":" + c.API_PORT + "/" + c.API_BASE;
	c.API_FILE_URL = c.API_URL + "/files";
	return c
}


