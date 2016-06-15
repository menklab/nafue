package nafue

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/menkveldj/nafue-api/models/display"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"github.com/menkveldj/nafue/models"
	"crypto/hmac"
	"crypto/sha256"
	"hash"
	"github.com/menkveldj/nafue/config"
)

var ()

func Decrypt(reader io.ReaderAt, writer io.WriterAt, password string, fileHeader *display.FileHeaderDisplay) (*models.FileBody, error) {

	//// get key
	//key := getPbkdf2(password, fileHeader.Salt)
	//
	//// decrypt
	//data, dErr := decrypt(secureData,nil, nil, key)
	//// if error decrypting return error
	//if dErr != nil {
	//	return &models.FileBody{}, dErr
	//}
	//
	//// use data to create a fileBody
	//var fileBody = models.FileBody{}
	//err := json.Unmarshal(*data, &fileBody)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return &fileBody, nil

	return nil, nil
}

func Encrypt(data io.ReaderAt, secureData io.ReadWriteSeeker, filename string, password string) (*display.FileHeaderDisplay, error) {

	// make salt
	salt, err := makeSalt();
	if err != nil {
		return nil, err
	}

	// make key with salt
	key := getPbkdf2(password, salt)
	block, err := aes.NewCipher(key)

	// make iv
	iv, err := makeIv(block.BlockSize())
	if err != nil {
		return nil, err
	}

	// create cipher stream
	stream := cipher.NewCTR(block, iv)

	// hash & checksum
	h := hmac.New(sha256.New, key)

	// do encryption
	_, err = encrypt(data, secureData, stream, h)
	if err != nil {
		return nil, err
	}

	// protect and hash filename
	sfn := []byte(config.FILENAME_KEY_START + filename + config.FILENAME_KEY_END)
	stream.XORKeyStream(sfn, sfn)
	//writer.Seek(0, 2)
	_, err = secureData.Write(sfn)
	if err != nil {
		return nil, err
	}
	h.Write(sfn)

	// append and hash iv
	_, err = secureData.Write(iv)
	if err != nil {
		return nil, err
	}
	h.Write(iv)

	// get hash sum
	mac := h.Sum(nil)
	_, err = secureData.Write(mac)
	if err != nil {
		return nil, err
	}

	// create file header
	fhd := display.FileHeaderDisplay{
		Salt: salt,
	}

	return &fhd, nil
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

func encrypt(data io.ReaderAt, secureData io.ReadWriteSeeker, stream cipher.Stream, h hash.Hash ) (int, error) {
	var i int = 0;
	var len int = 32000 // 32 kb
	var total int = 0;

	chunk := make([]byte, len)
	for {
		in, err := data.ReadAt(chunk, int64(len * i))
		if err != nil && err != io.EOF { // if its an error other than EOF
			return total, err
		} else if err == io.EOF && in == 0 { // if EOF and no extra data we are done
			return total, nil
		} else { // if extra data handle it before exit
			chunk = chunk[:in]
		}

		// do encryption
		stream.XORKeyStream(chunk, chunk)

		// write to file
		_, err = secureData.Write(chunk)
		if err != nil {
			return total, err
		}

		// add to hmac & md5
		out, _ := h.Write(chunk)

		total += out // add up all written io so we know how large file is
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
