package writing_bytes_to_file

import (
	"../writing_bytes_to_file"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestWriteBytestoFile(t *testing.T) {
	bytes, err := writing_bytes_to_file.WriteBytestoFile("test.txt", []byte("Ciao Andrea come stai?"))
	assert.Nil(t, err)
	assert.EqualValues(t, 22, bytes)
}

func TestQuickWriting(t *testing.T) {
	err := ioutil.WriteFile("test.txt", []byte("Ciao Andrea Come stai?"), 0666)
	assert.Nil(t, err)
}
