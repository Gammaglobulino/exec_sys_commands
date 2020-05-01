package hashing_generator

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestHashingFile(t *testing.T) {
	//open or create a file
	_, err := os.Stat("test.txt")
	var file *os.File
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File doesnt exist creating one")
			file, err := os.Create("test.txt")
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
		}
	}
	file.Close()

	data, err := ioutil.ReadFile("test.txt")
	MD5 := md5.Sum(data)
	SHA1 := sha1.Sum(data)
	SHA256 := sha256.Sum256(data)
	SHA512 := sha512.Sum512(data)
	assert.Nil(t, err)
	fmt.Printf("MD5: %s\n", base64.StdEncoding.EncodeToString(MD5[:]))
	fmt.Printf("SHA1: %s\n", base64.StdEncoding.EncodeToString(SHA1[:]))
	fmt.Printf("SHA256: %s\n", base64.StdEncoding.EncodeToString(SHA256[:]))
	fmt.Printf("SHA512: %s\n", base64.StdEncoding.EncodeToString(SHA512[:]))

}
