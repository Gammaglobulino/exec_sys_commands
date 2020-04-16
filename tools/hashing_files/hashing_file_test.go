package hashing_files

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
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
	assert.Nil(t, err)
	fmt.Printf("MD5: %x\n\n", md5.Sum(data))
	fmt.Printf("SHA1: %x\n\n", sha1.Sum(data))
	fmt.Printf("SHA256: %x\n\n", sha256.Sum256(data))
	fmt.Printf("SHA512: %x\n\n", sha512.Sum512(data))

}
