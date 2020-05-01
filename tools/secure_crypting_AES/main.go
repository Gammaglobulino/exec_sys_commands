package secure_crypting_AES

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
)

func GeneratePrimaryKey() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", nil
	}
	pk := base64.URLEncoding.EncodeToString(randomBytes)
	return pk, nil
}

func AESCryptFromToFile(key []byte, fromSourceFile string, toDestinationFileName string) error {
	//load the primary key
	readFileToEncrypt, err := ioutil.ReadFile(fromSourceFile)
	if err != nil {
		return err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	cipherTxt := make([]byte, aes.BlockSize+len(readFileToEncrypt))
	iv := cipherTxt[:aes.BlockSize]

	_, err = io.ReadFull(rand.Reader, iv) //generate random iv
	if err != nil {
		return err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherTxt[aes.BlockSize:], readFileToEncrypt)

	err = ioutil.WriteFile(toDestinationFileName, cipherTxt, 0666)
	if err != nil {
		return err
	}
	log.Println(string(readFileToEncrypt))
	log.Println(string(cipherTxt))
	return nil
}

func AESCryptFromDataBytesToByteArray(key []byte, bytesIn []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cipherTxt := make([]byte, aes.BlockSize+len(bytesIn))
	iv := cipherTxt[:aes.BlockSize]

	_, err = io.ReadFull(rand.Reader, iv) //generate random iv
	if err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherTxt[aes.BlockSize:], bytesIn)
	return cipherTxt, nil
}

func AESDecryptFromFile(key []byte, fromFileName string) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cipherText, err := ioutil.ReadFile(fromFileName)
	if err != nil {
		return nil, err
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(cipherText, cipherText)
	return cipherText, nil
}

func AESDecryptFromBytesToByteArray(key []byte, cipherText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(cipherText, cipherText)
	return cipherText, nil
}
