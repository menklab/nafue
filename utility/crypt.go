package utility

import (
	"golang.org/x/crypto/pbkdf2"
	"nafue/config"
	"io"
"crypto/rand"
	"log"
	"os"
)

//import "golang.org/x/crypto/pbkdf2"

var (

)

func Encrypt(password string){
	// create salt
	salt := makeSalt()
	log.Println("salt: ", salt)

	// generate key
	key := getPbkdf2(password, salt)
	log.Println("key: ", key)
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

func makeIV() {

}