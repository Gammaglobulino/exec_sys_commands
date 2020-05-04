package main

import (
	"../secure_crypting_AES"
	"../zip_archive"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/howeyc/gopass"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var text string
	var retries int
	pkey := "----"
	for text != pkey { // exit only if hash-verified passoword
		fmt.Print("Enter your password: ")
		password, err := gopass.GetPasswdMasked()
		if err != nil {
			log.Fatal(err)
		}
		enteredText := md5.Sum(password)
		text = base64.StdEncoding.EncodeToString(enteredText[:])
		if text != pkey {
			fmt.Println("Not a valid password, try again.")
			retries++
		}
		if retries == 3 {
			fmt.Println("You're an hacker, going to destroy the file!")
			os.Exit(1)
		}
	}

	cryptedData, err := zip_archive.UnzipFromFileToBytes("---")
	if err != nil {
		log.Fatal(err)
	}
	decryptedData, err := secure_crypting_AES.DecryptFromBytesToByteArray([]byte(pkey), cryptedData)
	if err != nil {
		fmt.Println("Corrupted data")
		os.Exit(1)
	}
	ioutil.WriteFile("mutuo.pdf", decryptedData, 0666)

}
