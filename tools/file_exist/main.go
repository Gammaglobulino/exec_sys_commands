package main

import (
	"log"
	"os"
)

func main() {
	fileInfo, err := os.Stat("test.txt")
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatal("File doesnt exist")
		}

	}
	log.Println(fileInfo)
}
