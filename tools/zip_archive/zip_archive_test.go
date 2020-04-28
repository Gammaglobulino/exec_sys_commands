package zip_archive

import (
	"../zip_archive"
	"archive/zip"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var filesToArchive = []struct{ Name, Data string }{
	{"test.txt", "String content for text 1"},
	{"test2.txt", "\x65\x66\x67"},
}

func TestZipBytesToFile(t *testing.T) {
	bytesToZip, err := ioutil.ReadFile("CV_Mazzanti Andrea.doc")
	if err != nil {
		t.Fatal(err)
	}
	nbytes, err := zip_archive.ZipBytesTo("test.zip", "AM.doc", bytesToZip)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Zipped %v bytes ", nbytes)
}
func TestUnzipBytesOriginalTo(t *testing.T) {
	err := zip_archive.UnzipBytesTo("", "test.zip")
	if err != nil {
		t.Fatal(err)
	}
}

func TestUnzipToDifferentFileName(t *testing.T) {
	err := zip_archive.UnzipBytesTo("AndreaMazzanti.doc", "test.zip")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddingandExtractingfilesToZipArchiveFiles(t *testing.T) {

	//Adding files to a zip Archive
	outFile, err := os.Create("test.zip")
	assert.Nil(t, err)
	defer outFile.Close()
	zipWriter := zip.NewWriter(outFile)
	for _, file := range filesToArchive {
		fileWriter, err := zipWriter.Create(file.Name)
		assert.Nil(t, err)
		_, err = fileWriter.Write([]byte(file.Data))
		assert.Nil(t, err)
	}
	err = zipWriter.Close()
	assert.Nil(t, err)

	t.Run("Extracting files from zipped archive", func(t *testing.T) {
		zipReader, err := zip.OpenReader("test.zip")
		assert.Nil(t, err)
		defer zipReader.Close()
		for _, file := range zipReader.Reader.File {
			zippedFile, err := file.Open()
			assert.Nil(t, err)
			defer zippedFile.Close()
			targetDir := "./" // current dir
			filePathToExtract := filepath.Join(targetDir, file.Name)
			if file.FileInfo().IsDir() {
				log.Println("Creating directory:", filePathToExtract)
				os.MkdirAll(filePathToExtract, file.Mode())
			} else {
				log.Println("Extracting file:", file.Name)
				outFile, err := os.OpenFile(filePathToExtract, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
				assert.Nil(t, err)
				defer outFile.Close()
				_, err = io.Copy(outFile, zippedFile)
				assert.Nil(t, err)
			}
		}
	})

}
