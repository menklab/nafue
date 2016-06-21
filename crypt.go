package nafue

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/menkveldj/nafue-api/models/display"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"crypto/hmac"
	"crypto/sha256"
	"github.com/menkveldj/nafue/config"
	"stash.cqlcorp.net/mp/moja-portal/utility/errors"
	"fmt"
	"os"
)

var C_DECRYPT_UNAUTHENTICATED error = errors.New("Data couldn't be authenticated. Is the password entered correct?")

func Decrypt(secureFile *os.File, password string, fileHeader *display.FileHeaderDisplay) (string, error) {

	//get key
	key := getPbkdf2(password, fileHeader.Salt)

	// start at start of file
	_, err := secureFile.Seek(0,0)
	if err != nil {
		return "", nil
	}

	// calculate mac1 from file
	h := hmac.New(sha256.New, key)
	fileSize, err := io.Copy(h, secureFile)
	if err != nil {
		return "", err
	}
	fileMac := h.Sum(nil)

	// verify hmac is good
	if ok := hmac.Equal(fileMac, fileHeader.Hmac); !ok {
		return "", C_DECRYPT_UNAUTHENTICATED
	}

	// get iv
	iv := make([]byte, aes.BlockSize, aes.BlockSize)
	fileSize -= aes.BlockSize // actual file doesn't contain iv
	_, err = secureFile.ReadAt(iv, fileSize)
	if err != nil {
		return "", err
	}

	// use iv to make cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// create cipher stream
	stream := cipher.NewCTR(block, iv)

	// decrypt
	err = decrypt(secureFile, stream, fileSize)
	if err != nil {
		return "", err
	}

	return "filename", nil
}

func Encrypt(file *os.File, secureFile *os.File, password string) (*display.FileHeaderDisplay, error) {

	// make salt
	salt, err := makeSalt();
	if err != nil {
		return nil, err
	}

	// make key with salt
	key := getPbkdf2(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// make iv
	iv, err := makeIv(block.BlockSize())
	if err != nil {
		return nil, err
	}

	// create cipher stream
	stream := cipher.NewCTR(block, iv)

	// do encryption
	err = encrypt(file, secureFile, stream)
	if err != nil {
		return nil, err
	}

	// protect and hash filename
	sfn := []byte(config.FILENAME_KEY_START + file.Name() + config.FILENAME_KEY_END)
	stream.XORKeyStream(sfn, sfn)
	_, err = secureFile.Write(sfn)
	if err != nil {
		return nil, err
	}

	// append iv
	_, err = secureFile.Write(iv)
	if err != nil {
		return nil, err
	}

	// get hash sum
	h := hmac.New(sha256.New, key)
	_, err = secureFile.Seek(0,0)	// start at start of file
	if err != nil {
		return nil, err
	}
	io.Copy(h, secureFile)
	mac := h.Sum(nil)

	// create file header
	fhd := display.FileHeaderDisplay{
		Salt: salt,
		Hmac: mac,
	}

	return &fhd, nil
}

func decrypt(secureData *os.File, stream cipher.Stream, eof int64) error {

	chunk := make([]byte, C.BUFFER_SIZE)
	// loop through file and decrypt
	for {
		fmt.Println("bytes left: ", eof)
		// if no more bytes; finish.
		if eof <= 0 {
			return nil
		}
		// if less bytes than chunk; shrink chunk
		if eof < C.BUFFER_SIZE {
			chunk = make([]byte, eof)
		}

		// read chunk
		readStart := eof - int64(len(chunk))
		_, err := secureData.ReadAt(chunk, readStart)
		if err != nil {
			return err
		}

		// decrypt
		fmt.Println("encrypted: ", chunk[:10])
		stream.XORKeyStream(chunk, chunk)
		fmt.Println("decrypted: ", chunk[:10])
		// write back to file
		if _, err = secureData.WriteAt(chunk, readStart); err != nil {
			return err
		}

		// reduce remaining bytes by amount read
		eof -= int64(len(chunk));
	}
	return nil
}

func encrypt(file *os.File, secureFile *os.File, stream cipher.Stream) error {
	var i int64 = 0;

	chunk := make([]byte, C.BUFFER_SIZE)
	for {
		in, err := file.ReadAt(chunk, int64(C.BUFFER_SIZE * i))
		if err != nil && err != io.EOF {
			// if its an error other than EOF
			return err
		} else if err == io.EOF && in == 0 {
			// if EOF and no extra data we are done
			return nil
		} else {
			// if extra data handle it before exit
			chunk = chunk[:in]
		}

		// do encryption
		stream.XORKeyStream(chunk, chunk)

		// write to file
		_, err = secureFile.WriteAt(chunk, int64(C.BUFFER_SIZE * i))
		if err != nil {
			return err
		}

		i++;
	}
	return nil
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
func makeIv(blockSize int) ([]byte, error) {
	iv := make([]byte, blockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	return iv, nil
}
