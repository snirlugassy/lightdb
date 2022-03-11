package object_db

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
	//Register(t interface{}) error
	Insert(object *interface{}) (int, error)
	Get(id int) interface{}
	Delete(id int)
	Update(id int, object *interface{}) error
	Commit() error
	Pull() error
}

//func (collection *Collection) Register(t interface{}) error {
//	if collection.Index != nil && len(collection.Index) == 0 {
//		collection.DType = reflect.TypeOf(t)
//	} else {
//		return errors.New("Type already registered for non-empty collection")
//	}
//	return nil
//}

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
	//data, marshelError := json.Marshal(collection.CollectionCore)
	//if marshelError != nil {
	//	return marshelError
	//}
	//writeFileError := ioutil.WriteFile(collection.FilePath + ".core", data, 0644)
	//if writeFileError != nil {
	//	return writeFileError
	//}

	/*	file, err := os.OpenFile(collection.FilePath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		// fmt.Println(reflect.New(collection.DType))
		// gob.Register(*collection)
		// gob.Register(reflect.New(collection.DType).Interface())
		gob.Register(collection.Index[1])
		encoder := gob.NewEncoder(file)
		err = encoder.Encode(collection)
		if err != nil {
			log.Fatal(err)
		}*/

	gob.RegisterName(reflect.New(collection.DType).String(), reflect.New(collection.DType).Interface())

	indexError := writeGob(collection.FilePath+".core", collection.CollectionCore)
	if indexError != nil {
		return indexError
	}

	//seqError := writeGob(collection.FilePath+".seq", collection.Seq)
	//if seqError != nil {
	//	return seqError
	//}

	//dtypeError := writeGob(collection.FilePath+".dtype", collection.DType)
	//if dtypeError != nil {
	//	return dtypeError
	//}
	//fmt.Println("D")

	return nil
}

func (collection *Collection) Pull() error {
	//data, readFileError := ioutil.ReadFile(collection.FilePath + ".core")
	//if readFileError != nil {
	//	return readFileError
	//}
	//
	//unmarshelError := json.Unmarshal(data, &collection.CollectionCore)
	//if unmarshelError != nil {
	//	log.Fatal(unmarshelError)
	//}
	//file, err := os.Open(collection.FilePath)
	//if err != nil {
	//	return err
	//}
	//defer file.Close()
	//e := new(interface{})
	//gob.Register(e)
	//decoder := gob.NewDecoder(file)
	//err = decoder.Decode(&collection.CollectionCore)
	//if err != nil {
	//	return err
	//}

	indexError := readGob(collection.FilePath+".core", &collection.CollectionCore)
	if indexError != nil {
		return indexError
	}

	//indexError := readGob(collection.FilePath+".index", &collection.Index)
	//if indexError != nil {
	//	return indexError
	//}
	//
	//seqError := readGob(collection.FilePath+".seq", &collection.Seq)
	//if seqError != nil {
	//	return seqError
	//}
	//
	////dtype := interface{}
	//dtypeError := readGob(collection.FilePath+".dtype", &collection.DType)
	//if dtypeError != nil {
	//	return dtypeError
	//}

	return nil
}

func writeGob(filePath string, object interface{}) error {
	file, writeFileError := os.Create(filePath)
	defer file.Close()
	if writeFileError != nil {
		return writeFileError
	}
	// gob.Register(object)
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
