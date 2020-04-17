package secure_crypting_AES

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"io/ioutil"
	"log"
)

func GeneratePrimaryKey() ([]byte, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, nil
	}
	return randomBytes, nil
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
