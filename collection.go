package lightdb

import (
	"encoding/gob"
	"errors"
	"reflect"
)

type CollectionCore struct {
	Index map[int]*interface{}
	Seq   int
}

type Collection struct {
	CollectionCore
	Name     string
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
	Find(query map[string]interface{}, results *[]interface{}) error
}

func (collection *Collection) Insert(object interface{}) (int, error) {
	if collection.Index == nil {
		collection.Index = make(map[int]*interface{})
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

func (collection *Collection) Update(id int, object interface{}) error {
	if reflect.TypeOf(object) != collection.DType {
		return errors.New("invalid type of object to update")
	}
	collection.Index[id] = &object
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

func (collection *Collection) Find(query map[string]interface{}, results *[]interface{}) error {
	validFields := make(map[string]interface{})

	for k, v := range query {
		if field, exists := collection.DType.FieldByName(k); exists && field.Type == reflect.TypeOf(v) {
			validFields[k] = v
		}
	}

	for _, item := range collection.Index {
		matchFlag := true
		for k, v := range validFields {
			if reflect.ValueOf(*item).FieldByName(k).Interface() != v {
				matchFlag = false
				break
			}
		}
		if matchFlag {
			*results = append(*results, *item)
		}
	}

	return nil
}
