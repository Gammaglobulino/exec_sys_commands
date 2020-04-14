package reading_from_file

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestReadingFromFile(t *testing.T) {

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
	assert.Nil(t, err)
	defer file.Close()

	bytesRed := make([]byte, 30)
	nbytes, err := file.Read(bytesRed)
	assert.Nil(t, err)
	assert.EqualValues(t, 30, nbytes)
	assert.EqualValues(t, "ABC <- scritti prima da bufio.", string(bytesRed))

	t.Run("Trying to fill the receiving buffer using io.ReadFull", func(t *testing.T) {
		file, err = os.Open("test.txt")
		assert.Nil(t, err)
		defer file.Close()
		bytesToRead := make([]byte, 15)
		nbytes, err := io.ReadFull(file, bytesToRead)
		assert.Nil(t, err)
		assert.EqualValues(t, 15, nbytes)
		assert.EqualValues(t, "ABC <- scritti ", string((bytesToRead)))

	})
	t.Run("Quick read the file data using ioutil.ReadFile", func(t *testing.T) {
		data, err := ioutil.ReadFile("test.txt")
		assert.Nil(t, err)
		assert.EqualValues(t, "ABC <- scritti prima da bufio.Write(byte[])", string(data))
	})

}
