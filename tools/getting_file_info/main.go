package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fileInfo, err := os.Stat("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File name:", fileInfo.Name())
	fmt.Println("File size:", fileInfo.Size())
	fmt.Println("File permissions:", fileInfo.Mode())
	fmt.Println("Last modified:", fileInfo.ModTime())
	fmt.Println("Is directory:", fileInfo.IsDir())
	fmt.Println("System Interface:", fileInfo.Sys())

}
