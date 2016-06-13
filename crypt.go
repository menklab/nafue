package nafue

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"github.com/menkveldj/nafue-api/models/display"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"github.com/menkveldj/nafue/models"
)

var ()

func Decrypt(fileHeader *display.FileHeaderDisplay, password string, secureData *[]byte) (*models.FileBody, error) {

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
	if err != nil {
		return nil, err
	}

	return &fileBody, nil
}

func Encrypt(reader io.Reader, key *[]byte, fileHeader *display.FileHeaderDisplay) (*[]byte, *display.FileHeaderDisplay, error) {


	//// marshal data for encryption
	//data, err := json.Marshal(*fileBodyPackage)
	//if err != nil {
	//	return nil, nil, err
	//}
	//
	//// encrypt
	//eData, err := encrypt(&data, aData, nonce, key)
	//return eData, fileDisplay, nil
	return nil, nil, nil
}

func decrypt(secureData *[]byte, aData []byte, nonce []byte, key []byte) (*[]byte, error) {

	//create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// decrypt
	gcm, err := cipher.NewGCM(block)
	data, err := gcm.Open(nil, nonce, *secureData, aData)
	if err != nil {
	}
	return &data, nil
}

func encrypt(data *[]byte, aData []byte, nonce []byte, key []byte) (*[]byte, error) {

	//create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// encrypt
	gcm, err := cipher.NewGCM(block)
	secureData := gcm.Seal(nil, nonce, *data, aData)

	return &secureData, nil
}

func getPbkdf2(password string, salt []byte) []byte {
	dk := pbkdf2.Key([]byte(password), salt, C.ITERATIONS, C.KEY_LENGTH, C.HASH_TYPE)
	return dk
}

func makeSalt() ([]byte, error) {
	salt := make([]byte, C.SALT_LENGTH)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}
func makeNonce() ([]byte, error) {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return nonce, nil
}

func makeAData() ([]byte, error) {
	aData := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, aData); err != nil {
		return nil, err
	}
	return aData, nil
}
