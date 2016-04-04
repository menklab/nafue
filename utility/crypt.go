package utility

import (
	"golang.org/x/crypto/pbkdf2"
	"nafue/config"
	"io"
	"crypto/rand"
	"log"
	"os"
	"errors"
	"io/ioutil"
	"encoding/base64"
	"crypto/aes"
	"crypto/cipher"
	"nafue/models"
	"encoding/json"
)


var ()

func Encrypt(file string, password string) {
	// create salt
	salt := makeSalt()
	log.Println("salt: ", salt)

	// generate key
	key := getPbkdf2(password, salt)
	log.Println("key: ", key)

	// encrypt
	encrypt(file, key)

}

func encrypt(file string, key []byte) {

	// verify file is under 50mb
	fileInfo, err := os.Stat(file)
	checkError(err)
	fileSize := fileInfo.Size()
	if fileSize > (config.FILE_SIZE_LIMIT * 1024 * 1024) {
		panic(errors.New("File is larger than " + string(config.FILE_SIZE_LIMIT) + "mb."))
	}

	// get file type and name
	fileName := fileInfo.Name()

	// read file
	fileBytes, err := ioutil.ReadFile(file)
	checkError(err)

	// converts bytes to base64
	b64 := base64.StdEncoding.EncodeToString(fileBytes)

	// create file data package
	fdp := models.FileDataPackage{
		Name: fileName,
		Content: b64,
	}
	data, err := json.Marshal(fdp)
	checkError(err)
	log.Println("Data: ", string(data))

	 //create cipher
	block, err := aes.NewCipher(key)

	// encrypt
	gcm, err := cipher.NewGCM(block)
	ciphertext := gcm.Seal(nil, makeNonce(), []byte(data), nil)
	log.Println("Cipher: ", ciphertext)

	// encode for storage
	encodedData := base64.StdEncoding.EncodeToString(data)
	log.Println("encodedData: ", encodedData)

}

func getPbkdf2(password string, salt []byte) []byte {
	dk := pbkdf2.Key([]byte(password), salt, config.ITERATIONS, config.KEY_LENGTH, config.HASH_TYPE)
	return dk
}

func decrypt() {

}

func makeSalt() []byte {
	salt := make([]byte, config.SALT_LENGTH)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		log.Println("Error creating salt: ", err.Error())
		os.Exit(1)
	}

	return salt
}
func makeNonce() []byte {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	return nonce
}