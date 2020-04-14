package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	//file doesnt exist? ok create
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
	file, err = os.Open("test.txt")

	fmt.Println(file.Name())
	defer file.Close()

	newFile, err := os.Create("new_file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()
	fmt.Println(newFile.Name())

	bytesWritten, err := io.Copy(newFile, file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Bytes written from %s to %s are %v ", file.Name(), newFile.Name(), bytesWritten)

}
