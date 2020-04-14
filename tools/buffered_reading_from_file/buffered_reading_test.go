package buffered_reading_from_file

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBufferedReading(t *testing.T) {
	file, err := os.OpenFile("scratch.txt", os.O_RDONLY, 0666)
	assert.Nil(t, err)
	assert.NotNil(t, file)
	defer file.Close()
	t.Run("Buffered peek the first 5 bytes", func(t *testing.T) {
		newreader := bufio.NewReader(file)
		bytesSlice, err := newreader.Peek(20)
		assert.Nil(t, err)
		assert.EqualValues(t, "Nel mezzo del cammin", string(bytesSlice))
	})
	t.Run("Read 5 bytes and advance the pointer", func(t *testing.T) {
		byteSlice5 := make([]byte, 5)
		newreader := bufio.NewReader(file)
		nbytes, err := newreader.Read(byteSlice5)
		assert.Nil(t, err)
		assert.EqualValues(t, 5, nbytes)
	})
	t.Run("Read string from the buffer", func(t *testing.T) {
		newreader := bufio.NewReader(file)
		data, err := newreader.ReadString('\n')
		assert.Nil(t, err)
		assert.EqualValues(t, "Nel mezzo del cammin di nostra vita\r\n", data)

	})
	t.Run("Reading with scanner", func(t *testing.T) {
		newscanner := bufio.NewScanner(file)
		newscanner.Split(bufio.ScanWords)
		ok := newscanner.Scan()
		assert.True(t, ok)
		assert.EqualValues(t, "Nel", newscanner.Text())
	})

}
