package downloading_file_over_http

import (
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestDownloadFromHttp(t *testing.T) {
	newFile, err := os.Create("test.html")
	assert.Nil(t, err)
	defer newFile.Close()

	url := "http://www.msn.com"
	response, err := http.Get(url)
	assert.Nil(t, err)
	defer response.Body.Close()
	bytesDownloaded, err := io.Copy(newFile, response.Body)
	assert.Nil(t, err)
	log.Println("Number of bytes downloaded", bytesDownloaded)

}
