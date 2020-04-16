package steganoziparchiveimage

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	image2 "image"
	"image/jpeg"
	"io"
	"log"
	"math/rand"
	"os"
	"testing"
)

//test files to be zipped
var filesToArchive = []struct{ Name, Data string }{
	{"test.txt", "String content for text 1"},
	{"test2.txt", "\x65\x66\x67"},
}

func TestCreateZipArchiveImage(t *testing.T) {

	//Create a zip Archive
	outZipFile, err := os.Create("test.zip")
	assert.Nil(t, err)
	zipWriter := zip.NewWriter(outZipFile)
	for _, file := range filesToArchive {
		fileWriter, err := zipWriter.Create(file.Name)
		assert.Nil(t, err)
		_, err = fileWriter.Write([]byte(file.Data))
		assert.Nil(t, err)
	}
	err = zipWriter.Close()
	assert.Nil(t, err)
	outZipFile.Close()

	//Creating an image

	image := image2.NewRGBA(image2.Rect(0, 0, 100, 200))
	for p := 0; p < (100 * 200); p++ {
		pixelOffset := 4 * p
		image.Pix[0+pixelOffset] = uint8(rand.Intn(256)) //R
		image.Pix[1+pixelOffset] = uint8(rand.Intn(256)) //G
		image.Pix[2+pixelOffset] = uint8(rand.Intn(256)) //B
		image.Pix[3+pixelOffset] = 255                   //Alpha
	}

	outJpgFile, err := os.Create("randomizedImage.jpg")
	assert.Nil(t, err)

	err = jpeg.Encode(outJpgFile, image, nil)
	assert.Nil(t, err)
	outJpgFile.Close()

	outZipFile, err = os.Open("test.zip")
	assert.Nil(t, err)
	defer outZipFile.Close()

	outJpgFile, err = os.Open("randomizedImage.jpg")
	assert.Nil(t, err)
	defer outJpgFile.Close()

	newFile, err := os.Create("stego_image.jpg")
	assert.Nil(t, err)
	defer newFile.Close()

	bytes, err := io.Copy(newFile, outJpgFile)
	assert.Nil(t, err)
	fmt.Println(bytes)

	bytes, err = io.Copy(newFile, outZipFile)
	assert.Nil(t, err)
	fmt.Println(bytes)

}
func TestDetectingStegoImageZipGene(t *testing.T) {
	filename := "stego_image.jpg"
	file, err := os.Open(filename)
	assert.Nil(t, err)
	bufferedReader := bufio.NewReader(file)
	fileStat, _ := file.Stat()
	var myByte byte
	byteSlice := make([]byte, 3)
	for i := int64(0); i < fileStat.Size(); i++ {
		myByte, err = bufferedReader.ReadByte()
		assert.Nil(t, err)
		if myByte == '\x50' {
			byteSlice, err = bufferedReader.Peek(3)
			assert.Nil(t, err)
			if bytes.Equal(byteSlice, []byte{'\x4b', '\x03', '\x04'}) {
				log.Println("Zip Archive found at byte:", i)
			}
		}
	}

}
