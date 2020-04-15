package compressing_files

import (
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestCompresstoGZIPfile(t *testing.T) {
	outGzipFile, err := os.Create("test.txt.gz")
	assert.Nil(t, err)

	gzipWriter := gzip.NewWriter(outGzipFile)
	_, err = gzipWriter.Write([]byte("Ciao sono Andrea Mazzanti\n"))
	assert.Nil(t, err)

	gzipWriter.Close()

	t.Run("Extracting gzipped file", func(t *testing.T) {
		gzipFile, err := os.Open("test.txt.gz")
		assert.Nil(t, err)

		gzipReader, err := gzip.NewReader(gzipFile)
		assert.Nil(t, err)
		defer gzipReader.Close()

		outFileWriter, err := os.Create("unzipped.txt")
		assert.Nil(t, err)

		defer outFileWriter.Close()

		_, err = io.Copy(outFileWriter, gzipReader)
		assert.Nil(t, err)
	})
}
