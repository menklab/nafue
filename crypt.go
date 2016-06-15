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
	"crypto/hmac"
	"crypto/sha256"
	"hash"
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

func Encrypt(reader io.ReaderAt, writer io.WriterAt, password string) (int64, error) {

	// make salt
	salt, err := makeSalt();
	if err != nil {
		return 0, err
	}

	// make key with salt
	key := getPbkdf2(password, salt)
	block, err := aes.NewCipher(key)

	// make iv
	iv, err := makeIv(block.BlockSize())
	if err != nil {
		return 0, err
	}

	// create cipher stream
	stream := cipher.NewCTR(block, iv)

	// hash
	h := hmac.New(sha256.New, key)

	// do encryption
	out, err := encrypt(reader, writer, stream, h)
	if err != nil {
		return 0, err
	}

	// append and hash iv
	writer.WriteAt(iv, out)
	h.Write(iv)

	// get hash sum
	mac := h.Sum(nil)
	writer.WriteAt(mac, (out + int64(len(iv))))

	return (out + int64(len(iv)) + int64(len(mac))), nil
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

func encrypt(reader io.ReaderAt, writer io.WriterAt, stream cipher.Stream, h hash.Hash) (int64, error) {
	var i int64 = 0;
	var len int64 = 32000 // 32 kb
	var total int64 = 0;

	data := make([]byte, len)
	for {
		in, err := reader.ReadAt(data, (len * i))
		if err != nil && err != io.EOF { // if its an error other than EOF
			return total, err
		} else if err == io.EOF && in == 0 { // if EOF and no extra data we are done
			return total, nil
		} else { // if extra data handle it before exit
			data = data[:in]
		}

		// do encryption
		stream.XORKeyStream(data, data)

		// write to file
		_, err = writer.WriteAt(data, (len*i))
		if err != nil {
			return total, err
		}

		// add to hmac
		out, err := h.Write(data)
		if err != nil {
			return total, err
		}
		total += int64(out) // add up all written io so we know how large file is
		i++;
	}
	return  total, nil
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
