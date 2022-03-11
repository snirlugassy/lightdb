package lightdb

import (
	"encoding/gob"
	"errors"
	"os"
	"reflect"
)

type CollectionCore struct {
	Index map[int]*interface{}
	Seq   int
}

type Collection struct {
	CollectionCore
	FilePath string
	DType    reflect.Type
}

type CollectionInterface interface {
	Insert(object *interface{}) (int, error)
	Get(id int) interface{}
	Delete(id int)
	Update(id int, object *interface{}) error
	Commit() error
	Pull() error
}

func (collection *Collection) Insert(object interface{}) (int, error) {
	if collection.Index == nil {
		collection.Index = make(map[int]*interface{})
		//collection.DType = reflect.TypeOf(object)
	} else if reflect.TypeOf(object) != collection.DType {
		return -1, errors.New("invalid data type for insert")
	}
	collection.Seq++
	collection.Index[collection.Seq] = &object
	return collection.Seq, nil
}

func (collection *Collection) Get(id int) interface{} {
	return *collection.CollectionCore.Index[id]
}

func (collection *Collection) Delete(id int) {
	delete(collection.Index, id)
}

func (collection *Collection) Update(id int, object *interface{}) error {
	if reflect.TypeOf(object) != collection.DType {
		return errors.New("invalid type of object to update")
	}
	collection.Index[id] = object
	return nil
}

func (collection *Collection) Commit() error {
	gob.RegisterName(reflect.New(collection.DType).String(), reflect.New(collection.DType).Interface())
	indexError := writeGob(collection.FilePath+".core", collection.CollectionCore)
	if indexError != nil {
		return indexError
	}
	return nil
}

func (collection *Collection) Pull() error {
	indexError := readGob(collection.FilePath+".core", &collection.CollectionCore)
	if indexError != nil {
		return indexError
	}
	return nil
}

func writeGob(filePath string, object interface{}) error {
	file, writeFileError := os.Create(filePath)
	defer file.Close()
	if writeFileError != nil {
		return writeFileError
	}
	encoder := gob.NewEncoder(file)
	encodingError := encoder.Encode(object)
	if encodingError != nil {
		return encodingError
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
