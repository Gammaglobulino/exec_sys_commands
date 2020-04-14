package main

import (
	"fmt"
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
	file.Close()

	file, err = os.OpenFile("test.txt", os.O_WRONLY, 0666)
	if err != nil {
		if os.IsPermission(err) {
			log.Println("Error: write permission denied")
		}
	}
	fmt.Println(file.Name())
	file.Close()

	err = os.Remove("test.txt")
	if err != nil {
		log.Fatal(err)
	}

}
