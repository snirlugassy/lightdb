package lightdb

import (
	"encoding/gob"
	"os"
)

func writeGob(filePath string, object interface{}) error {
	file, writeFileError := os.Create(filePath)
	if writeFileError != nil {
		return writeFileError
	}

	encoder := gob.NewEncoder(file)
	encodingError := encoder.Encode(object)
	if encodingError != nil {
		return encodingError
	}

	closeError := file.Close()
	if closeError != nil {
		return closeError
	}

	return nil
}

func readGob(filePath string, object interface{}) error {
	file, readFileError := os.Open(filePath)
	if readFileError != nil {
		return readFileError
	}
	decoder := gob.NewDecoder(file)
	decodingError := decoder.Decode(object)
	if decodingError != nil {
		return decodingError
	}
	file.Close()
	return nil
}
