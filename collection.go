package lightdb

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
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

func (collection *Collection) getCorePath() string {
	return collection.FilePath + ".core"
}

func (collection *Collection) Insert(object interface{}) (int, error) {
	collection.mutex.Lock()
	defer collection.mutex.Unlock()
	if collection.Index == nil {
		collection.Index = make(map[int]interface{})
	} else if reflect.TypeOf(object) != collection.DType {
		return -1, errors.New("invalid data type for insert")
	}
	collection.Seq++
	collection.Index[collection.Seq] = object
	return collection.Seq, nil
}

func (collection *Collection) InsertArray(objects []interface{}) ([]int, error) {
	collection.mutex.Lock()
	defer collection.mutex.Unlock()
	ids := make([]int, 0)
	for i := 0; i < len(objects); i++ {
		id, err := collection.Insert(objects[i])
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (collection *Collection) Get(id int) interface{} {
	return collection.CollectionCore.Index[id]
}

func (collection *Collection) GetAll() []interface{} {
	items := make([]interface{}, 0)
	for _, x := range collection.Index {
		items = append(items, x)
	}
	return items
}

func (collection *Collection) Delete(id int) {
	collection.mutex.Lock()
	defer collection.mutex.Unlock()
	delete(collection.Index, id)
}

func (collection *Collection) Update(id int, object interface{}) error {
	collection.mutex.Lock()
	defer collection.mutex.Unlock()
	if reflect.TypeOf(object) != collection.DType {
		return errors.New("invalid type of object to update")
	}
	collection.Index[id] = object
	return nil
}

func (collection *Collection) Commit() error {
	collection.mutex.Lock()
	defer collection.mutex.Unlock()
	gob.Register(reflect.New(collection.DType).Elem().Interface())
	err := writeObject(collection.getCorePath(), collection.CollectionCore)
	if err != nil {
		return err
	}
	return nil
}

func (collection *Collection) Pull() error {
	if _, err := os.Stat(collection.getCorePath()); os.IsNotExist(err) {
		return nil
	}

	collection.mutex.Lock()
	defer collection.mutex.Unlock()
	gob.Register(reflect.New(collection.DType).Elem().Interface())
	err := readObject(collection.FilePath+".core", &collection.CollectionCore)
	if err != nil {
		return err
	}
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

func (collection *Collection) ToJSON() ([]byte, error) {
	return json.Marshal(collection.GetAll())
}

func (collection *Collection) DirName() string {
	path, _ := os.Getwd()
	pathSlice := regexp.MustCompile("[^0-9A-Za-z_]").Split(path, -1)
	currentDir := pathSlice[len(pathSlice)-1]
	return currentDir

}

func (collection *Collection) CalculateFolderSize() (dirsize int64, err error) {

	err = os.Chdir(".")
	if err != nil {
		return
	}
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			dirsize += file.Size()
		}
	}
	return

}

func (collection *Collection) Info() (info map[string]interface{}) {
	folder := collection.DirName()

	fileSize, _ := collection.CalculateFolderSize()

	info = make(map[string]interface{})

	info["Folder"] = folder
	info["File Size(bytes)"] = fileSize
	info["collection name"] = collection.Name
	return

}
