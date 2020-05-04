package secure_crypting_AES

import (
	"../secure_crypting_AES"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
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
	pk := base64.URLEncoding.EncodeToString(randomBytes)
	log.Println(pk)

}

func TestStoreKeytoFileAndVerify(t *testing.T) {
	filename := "key.bin"
	pk, err := secure_crypting_AES.GeneratePrimaryKey()
	assert.Nil(t, err)
	err = ioutil.WriteFile(filename, []byte(pk), 0666)
	assert.Nil(t, err)
	readData, err := ioutil.ReadFile("key.bin")
	assert.Nil(t, err)
	assert.EqualValues(t, pk, readData)
}

func TestAESCrypt(t *testing.T) {
	//load the primary key
	readFileToEncrypt, err := ioutil.ReadFile("mutuo_peschiera_casa.pdf")
	if err != nil {
		t.Fatal(err)
	}

	readPK, err := ioutil.ReadFile("amkey.key")
	if err != nil {
		t.Fatal(err)
	}

	block, err := aes.NewCipher(readPK)
	if err != nil {
		t.Fatal(err)
	}

	cipherTxt := make([]byte, aes.BlockSize+len(readFileToEncrypt))

	iv := cipherTxt[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv) //generate random iv
	assert.Nil(t, err)
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherTxt[aes.BlockSize:], readFileToEncrypt)

	err = ioutil.WriteFile("encripted_mutuo.dat", cipherTxt, 0666)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAESDecryptFromBytes(t *testing.T) {
	key := []byte("kdjfheyrnchdtescdetfbskdiehdtass")
	data := []byte("Testo di prova per capire se va")
	cryptedData, err := secure_crypting_AES.CryptBytesArrayToByteArray(key, data)
	if err != nil {
		t.Fatal(err)
	}
	originalData, err := secure_crypting_AES.DecryptFromBytesToByteArray(key, cryptedData)
	if err != nil {
		t.Fatal()
	}
	assert.EqualValues(t, data, originalData)

}

func TestDecryptAESCryptedFile(t *testing.T) {

	readPK, err := ioutil.ReadFile("amkey.key")
	assert.Nil(t, err)

	block, err := aes.NewCipher(readPK)
	assert.Nil(t, err)

	cipherText, err := ioutil.ReadFile("encripted_mutuo.dat")
	assert.Nil(t, err)

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(cipherText, cipherText)
	ioutil.WriteFile("mutuo.pdf", cipherText, 0666)
}

func TestAESCryptFromToFile(t *testing.T) {
	readPK, err := ioutil.ReadFile("key.bin")
	assert.Nil(t, err)
	err = secure_crypting_AES.CryptFileToFile(readPK, "commedia_canto_primo", "comedycrypt.dat")
	assert.Nil(t, err)
}
func TestAESDecryptFromFile(t *testing.T) {
	readPK, err := ioutil.ReadFile("key.bin")
	assert.Nil(t, err)
	data, err := secure_crypting_AES.DecryptFromFile(readPK, "comedycrypt.dat")
	assert.Nil(t, err)
	assert.NotEmpty(t, data)
	log.Println(string(data))
}
