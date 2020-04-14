package buffered_writing_to_file

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestBufferedBytesToFile(t *testing.T) {

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
	assert.Nil(t, err)
	assert.NotNil(t, file)
	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)

	t.Run("Test buffer space available before writing must be 4096", func(t *testing.T) {
		bytesAvailable := bufferedWriter.Available()
		assert.EqualValues(t, 4096, bytesAvailable)
	})

	bytes, err := bufferedWriter.Write([]byte{65, 66, 67})
	assert.Nil(t, err)
	assert.EqualValues(t, 3, bytes)

	bytes, err = bufferedWriter.WriteString(" <- scritti prima da bufio.Write(byte[])")
	assert.Nil(t, err)
	assert.EqualValues(t, 40, bytes)

	t.Run("Test buffer space available after writing must be 4053", func(t *testing.T) {
		bytesAvailable := bufferedWriter.Available()
		assert.EqualValues(t, 4053, bytesAvailable)
	})

	t.Run("Unflushed bytes waiting must be 43", func(t *testing.T) {
		bytesWaiting := bufferedWriter.Buffered()
		assert.EqualValues(t, 43, bytesWaiting)

	})

	bufferedWriter.Flush()

}
