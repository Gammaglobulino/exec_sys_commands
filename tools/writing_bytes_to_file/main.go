package writing_bytes_to_file

import (
	"os"
)

func WriteBytestoFile(name string, bytes []byte) (int, error) {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	bytesWritten, err := file.Write(bytes)
	if err != nil {
		return 0, err
	}
	return bytesWritten, nil

}
