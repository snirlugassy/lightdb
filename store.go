package lightdb

import (
	"encoding/gob"
	"os"
)

func writeObject(filePath string, object interface{}) error {
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

func readObject(filePath string, object interface{}) error {
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

// func writeJSON(filePath string, object interface{}) error {
// 	data, encodingError := json.Marshal(object)
// 	if encodingError != nil {
// 		return encodingError
// 	}

// 	file, writeFileError := os.Create(filePath)
// 	if writeFileError != nil {
// 		return writeFileError
// 	}

// 	file.Write(data)
// 	file.Close()
// 	return nil
// }

// func readJSON(filePath string, target interface{}) error {
// 	data, err := os.ReadFile(filePath)
// 	if err != nil {
// 		return err
// 	}

// 	err = json.Unmarshal(data, target)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
