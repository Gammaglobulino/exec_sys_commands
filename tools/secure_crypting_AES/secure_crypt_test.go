package secure_crypting_AES

import (
	"../secure_crypting_AES"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"log"
	"testing"
)

func TestGeneratePrimaryKey(t *testing.T) {
	randomBytes := make([]byte, 32)
	numBytesRead, err := rand.Read(randomBytes)
	assert.Nil(t, err)
	assert.EqualValues(t, 32, numBytesRead)
	log.Println(randomBytes)
}

func TestStoreKeytoFileAndVerify(t *testing.T) {
	filename := "key.bin"
	pk, err := secure_crypting_AES.GeneratePrimaryKey()
	assert.Nil(t, err)
	err = ioutil.WriteFile(filename, pk, 0666)
	assert.Nil(t, err)
	readData, err := ioutil.ReadFile("key.bin")
	assert.Nil(t, err)
	assert.EqualValues(t, pk, readData)
}

func TestAESCrypt(t *testing.T) {
	//load the primary key
	readFileToEncrypt, err := ioutil.ReadFile("commedia_canto_primo")
	assert.Nil(t, err)

	readPK, err := ioutil.ReadFile("key.bin")
	assert.Nil(t, err)

	block, err := aes.NewCipher(readPK)
	assert.Nil(t, err)

	cipherTxt := make([]byte, aes.BlockSize+len(readFileToEncrypt))

	iv := cipherTxt[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv) //generate random iv
	assert.Nil(t, err)
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherTxt[aes.BlockSize:], readFileToEncrypt)

	err = ioutil.WriteFile("encripted_comedy.dat", cipherTxt, 0666)
	assert.Nil(t, err)
	log.Println(string(readFileToEncrypt))
	log.Println(string(cipherTxt))
}

func TestDecryptAESCryptedFile(t *testing.T) {

	readPK, err := ioutil.ReadFile("key.bin")
	assert.Nil(t, err)

	block, err := aes.NewCipher(readPK)
	assert.Nil(t, err)

	cipherText, err := ioutil.ReadFile("comedycrypt.dat")
	assert.Nil(t, err)

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(cipherText, cipherText)
	log.Println(string(cipherText))
}

func TestAESCryptFromToFile(t *testing.T) {
	readPK, err := ioutil.ReadFile("key.bin")
	assert.Nil(t, err)
	err = secure_crypting_AES.AESCryptFromToFile(readPK, "commedia_canto_primo", "comedycrypt.dat")
	assert.Nil(t, err)
}
func TestAESDecryptFromFile(t *testing.T) {
	readPK, err := ioutil.ReadFile("key.bin")
	assert.Nil(t, err)
	data, err := secure_crypting_AES.AESDecryptFromFile(readPK, "comedycrypt.dat")
	assert.Nil(t, err)
	assert.NotEmpty(t, data)
	log.Println(string(data))

}
