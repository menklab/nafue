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
)

var ()

func Decrypt(secureData io.ReadWriteSeeker, password string, fileHeader *display.FileHeaderDisplay) (error) {

	//get key
	key := getPbkdf2(password, fileHeader.Salt)

	// calculate mac1 from file
	h := hmac.New(sha256.New, key)
	_, err := io.Copy(h, secureData)
	if err != nil {
		return err
	}
	fileMac := h.Sum(nil)

	// get mac2 at end of file
	_, err = secureData.Seek(-sha256.Size, 2)
	if err != nil {
		return err
	}

	// verify hmac is good
	if ok := hmac.Equal(fileMac, fileHeader.Hmac); !ok {
		return errors.New("Data couldn't be authenticated. Is the password entered correct?")
	}

	// get iv
	iv := make([]byte, aes.BlockSize, aes.BlockSize)
	eof, err := secureData.Seek(-aes.BlockSize, 2)
	if err != nil {
		return err
	}
	_, err = secureData.Read(iv)
	if err != nil {
		return err
	}

	// use iv to make cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// create cipher stream
	stream := cipher.NewCTR(block, iv)

	// decrypt
	err = decrypt(secureData, stream, eof)
	if err != nil {
		return err
	}

	// use data to create a fileBody
	//var fileBody = models.FileBody{}
	//err := json.Unmarshal(*data, &fileBody)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return &fileBody, nil

	return nil
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
	err = encrypt(data, secureData, stream)
	if err != nil {
		return nil, err
	}

	// protect and hash filename
	sfn := []byte(config.FILENAME_KEY_START + filename + config.FILENAME_KEY_END)
	stream.XORKeyStream(sfn, sfn)
	_, err = secureData.Write(sfn)
	if err != nil {
		return nil, err
	}

	// append iv
	_, err = secureData.Write(iv)
	if err != nil {
		return nil, err
	}

	// get hash sum
	// hash & checksum
	h := hmac.New(sha256.New, key)
	io.Copy(h, secureData)
	mac := h.Sum(nil)

	// create file header
	fhd := display.FileHeaderDisplay{
		Salt: salt,
		Hmac: mac,
	}

	return &fhd, nil
}

func decrypt(secureData io.ReadWriteSeeker, stream cipher.Stream, eof int64) error {
	var length int64 = 32000 // 32 kb

	// start of file
	if _, err := secureData.Seek(0,0); err != nil {
		return err
	}
	chunk := make([]byte, length)
	// loop through file and decrypt
	for {
		fmt.Println("bytes left: ", eof)
		// if no more bytes; finish.
		if eof <= 0 {
			return nil
		}
		// if less bytes than chunk; shrink chunk
		if eof < length {
			chunk = make([]byte, eof)
		}

		// read chunk
		in, err := secureData.Read(chunk)
		if err != nil {
			return err
		}

		// decrypt
		fmt.Println("encrypted: ", chunk[:10])
		stream.XORKeyStream(chunk, chunk)
		fmt.Println("decrypted: ", chunk[:10])
		// write back to file
		if _, err = secureData.Seek(int64(-in),1); err != nil {
			return err
		}
		if _, err = secureData.Write(chunk); err != nil {
			return err
		}

		// reduce remaining bytes by amount read
		eof -= int64(in);
	}
	return nil
}

func encrypt(data io.ReaderAt, secureData io.ReadWriteSeeker, stream cipher.Stream) error {
	var i int = 0;
	var len int = 32000 // 32 kb

	chunk := make([]byte, len)
	for {
		in, err := data.ReadAt(chunk, int64(len * i))
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
		_, err = secureData.Write(chunk)
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
