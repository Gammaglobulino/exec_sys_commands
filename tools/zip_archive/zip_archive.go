package zip_archive

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
)

func ZipBytesTo(outZipFileName string, destinationFileName string, from []byte) (int, error) {

	outFile, err := os.Create(outZipFileName)
	if err != nil {
		return 0, err
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)

	fileWriter, err := zipWriter.Create(destinationFileName)
	nbytes, err := fileWriter.Write(from)
	if err != nil {
		return 0, nil
	}
	err = zipWriter.Close()
	if err != nil {
		return 0, nil
	}
	return nbytes, nil

}

func UnzipBytesTo(destinationFilename string, fromZipFile string) error {
	zipReader, err := zip.OpenReader(fromZipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, file := range zipReader.Reader.File {
		zippedFile, err := file.Open()
		if err != nil {
			return err
		}
		defer zippedFile.Close()
		targetDir := "./" // current dir
		if destinationFilename == "" {
			destinationFilename = file.FileInfo().Name()
		}
		filePathToExtract := filepath.Join(targetDir, destinationFilename)
		if file.FileInfo().IsDir() {
			log.Println("Creating directory:", filePathToExtract)
			os.MkdirAll(filePathToExtract, file.Mode())
		} else {
			log.Println("Extracting file:", destinationFilename)
			outFile, err := os.OpenFile(filePathToExtract, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, zippedFile)
			if err != nil {
				return err
			}

		}
	}
	return nil
}
