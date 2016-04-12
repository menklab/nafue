package utility

import (
	"golang.org/x/crypto/pbkdf2"
	"nafue/config"
	"io"
	"crypto/rand"
	"log"
	"os"
	"crypto/aes"
	"crypto/cipher"
	"nafue/models"
	"encoding/json"
	"nafue-api/models/display"
	"errors"
)


var ()

func Decrypt (fileHeader *display.FileHeaderDisplay, password string, secureData *[]byte) (*models.FileBody, error) {

	// get key
	key := getPbkdf2(password, fileHeader.Salt)

	// decrypt
	data, dErr := decrypt(secureData, fileHeader.AData, fileHeader.IV, key)
	// if error decrypting return error
	if dErr != nil {
		return &models.FileBody{}, dErr
	}

	// use data to create a fileBody
	var fileBody = models.FileBody{}
	err := json.Unmarshal(*data, &fileBody)
	checkError(err)

	return &fileBody, nil
}

func Encrypt(fileBodyPackage *models.FileBody, password string) (*[]byte, *display.FileHeaderDisplay) {
	// create aData
	aData := makeAData()

	// create salt
	salt := makeSalt()

	// create nonce
	nonce := makeNonce()

	// generate key
	key := getPbkdf2(password, salt)

	// create file display
	fileDisplay := display.FileHeaderDisplay{
		Salt: salt,
		// Todo update IV to nonce once api and ui is updated
		IV: nonce,
		AData: aData,
	}

	// marshal data for encryption
	data, err := json.Marshal(*fileBodyPackage)
	checkError(err)

	// encrypt
	return encrypt(&data, aData, nonce, key), &fileDisplay
}

func decrypt(secureData *[]byte, aData []byte, nonce []byte, key []byte) (*[]byte, error) {

	//create cipher
	block, err := aes.NewCipher(key)
	checkError(err)

	// decrypt
	gcm, err := cipher.NewGCM(block)
	data, err := gcm.Open(nil, nonce, *secureData, aData)
	if err != nil {
		return &data, errors.New("Bad Password")
	}
	return &data, nil
}

func encrypt(data *[]byte, aData []byte, nonce []byte, key []byte) *[]byte{

	 //create cipher
	block, err := aes.NewCipher(key)
	checkError(err)

	// encrypt
	gcm, err := cipher.NewGCM(block)
	secureData := gcm.Seal(nil, nonce, *data, aData)

	return &secureData
}

func getPbkdf2(password string, salt []byte) []byte {
	dk := pbkdf2.Key([]byte(password), salt, config.ITERATIONS, config.KEY_LENGTH, config.HASH_TYPE)
	return dk
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

func makeAData() []byte {
	aData := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, aData); err != nil {
		panic(err.Error())
	}
	return aData
}