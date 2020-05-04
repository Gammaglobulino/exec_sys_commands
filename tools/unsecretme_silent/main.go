package main

import (
	"../secure_crypting_AES"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	//parameters 1st : password 2nd sourceFilename 3rd:destinationFilename

	argsWithoutProg := os.Args[1:]
	pkey := argsWithoutProg[0]
	sourceFilename := argsWithoutProg[1]
	destinationFilename := argsWithoutProg[2]
	if pkey == "" {
		os.Exit(1)
	}

	_, err := os.Stat(sourceFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cryptedData, err := ioutil.ReadFile(sourceFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	decryptedData, err := secure_crypting_AES.DecryptFromBytesToByteArray([]byte(pkey), cryptedData)
	if err != nil {
		fmt.Println("Corrupted data")
		os.Exit(1)
	}
	ioutil.WriteFile(destinationFilename, decryptedData, 0666)
}
