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
	"nafue-api/models/display"
)


var ()

func Decrypt (fileDisplay *display.FileDisplay, password string, secureData *string) *string{

	// get aData
	aData, err := base64.StdEncoding.DecodeString(fileDisplay.AData)
	checkError(err)
	log.Println("aData: ", aData)

	// get salt
	salt, err := base64.StdEncoding.DecodeString(fileDisplay.Salt)
	checkError(err)
	log.Println("salt: ", salt)

	// get nonce
	nonce, err := base64.StdEncoding.DecodeString(fileDisplay.IV)
	checkError(err)
	log.Println("salt: ", nonce)

	// get key
	key := getPbkdf2(password, salt)
	log.Println("key: ", key)

	// decrypt
	decrypt(secureData, aData, nonce, key)

	return nil
}

func Encrypt(file string, password string) (*string, *display.FileDisplay) {
	// create aData
	aData := makeAData()
	log.Println("aData: ", aData)

	// create salt
	salt := makeSalt()
	log.Println("salt: ", salt)

	// create nonce
	nonce := makeNonce()
	log.Println("nonce: ", nonce)

	// generate key
	key := getPbkdf2(password, salt)
	log.Println("key: ", key)

	// encrypt
	return encrypt(file, aData, nonce, key), &display.FileDisplay{
		Salt: base64.StdEncoding.EncodeToString(salt),
		// Todo update IV to nonce once api and ui is updated
		IV: base64.StdEncoding.EncodeToString(nonce),
		AData: base64.StdEncoding.EncodeToString(aData),
	}



}

func decrypt(secureData *string, aData []byte, nonce []byte, key []byte) *string {

	// decode from storage
	decodedData, err := base64.StdEncoding.DecodeString(*secureData)
	checkError(err)
	log.Println("decoded: ", decodedData)

	// decrypt data
	//create cipher
	block, err := aes.NewCipher(key)

	// decrypt
	gcm, err := cipher.NewGCM(block)
	decryptedData, err := gcm.Open(nil, nonce, []byte(*secureData), aData)
	checkError(err)
	log.Println("Decrypted Data: ", decryptedData)

	return nil
}

func encrypt(file string, aData []byte, nonce []byte, key []byte) *string{

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

	// create file data package
	fdp := models.FileDataPackage{
		Name: fileName,
		Content: fileBytes,
	}
	// todo probably need to encode this data for it to work properl
	data, err := json.Marshal(fdp)
	checkError(err)
	log.Println("Data: ", string(data))

	 //create cipher
	block, err := aes.NewCipher(key)

	// encrypt
	gcm, err := cipher.NewGCM(block)
	ciphertext := gcm.Seal(nil, nonce, []byte(data), aData)
	log.Println("Cipher: ", ciphertext)

	// encode for storage
	encodedData := base64.StdEncoding.EncodeToString(data)
	log.Println("encodedData: ", encodedData)
	return &encodedData

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
	aData := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, aData); err != nil {
		panic(err.Error())
	}
	return aData
}