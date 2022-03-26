package lightdb

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"reflect"
)

type Database struct {
	Collections []*Collection
	Name        string
	Path        string
}

func (db *Database) Init() error {
	log.Println(db.Path)
	db.Collections = make([]*Collection, 0)
	_, err := os.Open(db.Path)

	if os.IsNotExist(err) {
		// If the path doesn't exists, try to mkdir
		mkdirError := os.Mkdir(db.Path, 0666)
		if mkdirError != nil {
			return mkdirError
		}
	}

	pathStat, err := os.Stat(db.Path)
	if err != nil {
		return err
	}

	if os.IsExist(err) && !pathStat.IsDir() {
		// If the path exists but isn't a dir -> raise error
		return errors.New("db path exists but not a folder")
	}

	return nil
}

func (db *Database) InitCollection(name string, dtype reflect.Type) *Collection {
	collection := Collection{
		Name:     name,
		DType:    dtype,
		FilePath: filepath.Join(db.Path, name+".db"),
	}
	db.Collections = append(db.Collections, &collection)
	return &collection
}

func (db *Database) LoadCollection(name string, dtype reflect.Type) (*Collection, error) {
	collection := db.InitCollection(name, dtype)
	pullError := collection.Pull()
	db.Collections = append(db.Collections, collection)
	return collection, pullError
}

func (db *Database) GetCollection(name string) (*Collection, bool) {
	for i := 0; i < len(db.Collections); i++ {
		if db.Collections[i].Name == name {
			return db.Collections[i], true
		}
	}
	return &Collection{}, false
}

func (db *Database) Commit() error {
	for i := 0; i < len(db.Collections); i++ {
		err := db.Collections[i].Commit()
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *Database) Pull() error {
	for i := 0; i < len(db.Collections); i++ {
		err := db.Collections[i].Pull()
		if err != nil {
			return err
		}
	}
	return nil
}
