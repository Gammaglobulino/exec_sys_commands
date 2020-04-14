package main

import (
	"log"
	"os"
)

func main() {
	newFile, err := os.Create("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()
	log.Println(newFile.Name())
}
