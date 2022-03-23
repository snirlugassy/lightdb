package lightdb

import (
	"encoding/gob"
	"errors"
	"reflect"
	"sync"
)

type CollectionCore struct {
	Index map[int]interface{}
	Seq   int
}

type Collection struct {
	CollectionCore
	Name     string
	FilePath string
	DType    reflect.Type
	mutex    sync.Mutex
}

//type CollectionInterface interface {
//	Insert(object interface{}) (int, error)
//	InsertArray(objects []interface{}) ([]int, error)
//	Get(id int) interface{}
//	Delete(id int)
//	Update(id int, object interface{}) error
//	Commit() error
//	Pull() error
//	Find(query map[string]interface{}, results *[]interface{})
//	First(query map[string]interface{}, result *interface{}) error
//	Filter(query map[string]interface{}, results *[]interface{})
//	FilterFirst(query map[string]interface{}, result *interface{})
//}

func (collection *Collection) Insert(object interface{}) (int, error) {
	collection.mutex.Lock()
	if collection.Index == nil {
		collection.Index = make(map[int]interface{})
	} else if reflect.TypeOf(object) != collection.DType {
		return -1, errors.New("invalid data type for insert")
	}
	collection.Seq++
	collection.Index[collection.Seq] = object
	collection.mutex.Unlock()
	return collection.Seq, nil
}

func (collection *Collection) InsertArray(objects []interface{}) ([]int, error) {
	collection.mutex.Lock()
	ids := make([]int, 0)
	for i := 0; i < len(objects); i++ {
		id, err := collection.Insert(objects[i])
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	collection.mutex.Unlock()
	return ids, nil
}

func (collection *Collection) Get(id int) interface{} {
	return collection.CollectionCore.Index[id]
}

func (collection *Collection) Delete(id int) {
	collection.mutex.Lock()
	delete(collection.Index, id)
	collection.mutex.Unlock()
}

func (collection *Collection) Update(id int, object interface{}) error {
	collection.mutex.Lock()
	if reflect.TypeOf(object) != collection.DType {
		return errors.New("invalid type of object to update")
	}
	collection.Index[id] = object
	collection.mutex.Unlock()
	return nil
}

func (collection *Collection) Commit() error {
	collection.mutex.Lock()
	gob.Register(reflect.New(collection.DType).Elem().Interface())
	err := writeObject(collection.FilePath+".core", collection.CollectionCore)
	if err != nil {
		return err
	}
	collection.mutex.Unlock()
	return nil
}

func (collection *Collection) Pull() error {
	collection.mutex.Lock()
	gob.Register(reflect.New(collection.DType).Elem().Interface())
	err := readObject(collection.FilePath+".core", &collection.CollectionCore)
	if err != nil {
		return err
	}
	collection.mutex.Unlock()
	return nil
}

func (collection *Collection) Find(query map[string]interface{}, results *[]interface{}) {
	validFields := make(map[string]interface{})

	for k, v := range query {
		if field, exists := collection.DType.FieldByName(k); exists && field.Type == reflect.TypeOf(v) {
			validFields[k] = v
		}
	}

	for _, item := range collection.Index {
		matchFlag := true
		for k, v := range validFields {
			if reflect.ValueOf(item).FieldByName(k).Interface() != v {
				matchFlag = false
				break
			}
		}
		if matchFlag {
			*results = append(*results, item)
		}
	}
}

func (collection *Collection) First(query map[string]interface{}, result *interface{}) {
	validField := make(map[string]interface{})

	for k, v := range query {
		if field, exists := collection.DType.FieldByName(k); exists && field.Type == reflect.TypeOf(v) {
			validField[k] = v
		}
	}

	for _, item := range collection.Index {
		matchFlag := true
		for k, v := range validField {
			if reflect.ValueOf(item).FieldByName(k).Interface() != v {
				matchFlag = false
				break
			}
		}
		if matchFlag {
			*result = item
			break
		}
	}
}

func (collection *Collection) Filter(filter func(v interface{}) bool, results *[]interface{}) {
	for _, item := range collection.Index {
		if filter(item) {
			*results = append(*results, item)
		}
	}
}

func (collection *Collection) FilterFirst(filter func(v interface{}) bool, result *interface{}) {
	for _, item := range collection.Index {
		if filter(item) {
			*result = item
			break
		}
	}
}
