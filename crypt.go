package nafue

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/menkveldj/nafue-api/models"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"crypto/hmac"
	"crypto/sha256"
	"os"
	"path"
	"github.com/menkveldj/nafue-api/utility/errors"
)

var C_DECRYPT_UNAUTHENTICATED error = errors.New("Data couldn't be authenticated. Is the password entered correct?")
var padding byte = []byte("!")[0]

func Decrypt(secureFile *os.File, password string, fileHeader *models.FileHeader) (string, error) {

	//get key
	key := getPbkdf2(password, fileHeader.Salt)

	// start at start of file
	_, err := secureFile.Seek(0, 0)
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
	err = XORKeyStreamBuffered(secureFile, secureFile, stream)
	if err != nil {
		return "", err
	}

	//truncate iv
	if err = os.Truncate(secureFile.Name(), fileSize); err != nil {
		return "", err
	}

	// get filename
	//read last C.MAX_FILENAME_LENGTH from file
	paddedFileName := make([]byte, C.MAX_FILENAME_LENGTH)
	_, err = secureFile.ReadAt(paddedFileName, (fileSize - C.MAX_FILENAME_LENGTH))
	if err != nil {
		return "", err
	}
	eofn := firstIndex(paddedFileName, make([]byte, 1)[0])
	if eofn == -1 {
		return "", errors.New("Bad filename enclosed in file.")
	}
	fileName := string(paddedFileName[:eofn])
	fileSize -= C.MAX_FILENAME_LENGTH

	// truncate filename so file is as original
	if err = os.Truncate(secureFile.Name(), fileSize); err != nil {
		return "", err
	}

	return fileName, nil
}

func Encrypt(file *os.File, secureFile *os.File, password string) (*models.FileHeader, error) {

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
	err = XORKeyStreamBuffered(file, secureFile, stream)
	if err != nil {
		return nil, err
	}

	// seek to end
	if _, err = secureFile.Seek(0, 2); err != nil {
		return nil, err
	}

	// protect filename
	paddedFileName := make([]byte, C.MAX_FILENAME_LENGTH)
	fileName := []byte(path.Base(file.Name()))
	if int64(len(fileName)) > C.MAX_FILENAME_LENGTH {
		return nil, errors.New("Filename must be less than or equal to 225 chars.")
	}
	copy(paddedFileName, fileName)
	stream.XORKeyStream(paddedFileName, paddedFileName)

	_, err = secureFile.Write(paddedFileName)
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
	_, err = secureFile.Seek(0, 0)        // start at start of file
	if err != nil {
		return nil, err
	}
	io.Copy(h, secureFile)
	mac := h.Sum(nil)

	// create file header
	fhd := models.FileHeader{
		Salt: salt,
		Hmac: mac,
	}

	return &fhd, nil
}

func XORKeyStreamBuffered(file *os.File, secureFile *os.File, stream cipher.Stream) error {
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

func firstIndex(s []byte, e byte) int {
	for i, _ := range s {
		if s[i] == e {
			return i
		}
	}
	return -1
}
