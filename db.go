package lightdb

import (
	"errors"
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
	db.Collections = make([]*Collection, 0)
	dbPath, err := os.Open(db.Path)
	if err != nil {
		return err
	}

	pathInfo, err := dbPath.Stat()
	if err != nil {
		return err
	}

	if !pathInfo.IsDir() {
		errors.New("db path is not a directory")
	}

	return nil
}

func (db *Database) CreateCollection(name string, dtype reflect.Type) *Collection {
	collection := Collection{
		Name:     name,
		DType:    dtype,
		FilePath: filepath.Join(db.Path, name+".db"),
	}
	db.Collections = append(db.Collections, &collection)
	return &collection
}

func (db *Database) LoadCollection(name string, dtype reflect.Type) (*Collection, error) {
	collection := Collection{
		Name:     name,
		DType:    dtype,
		FilePath: filepath.Join(db.Path, name+".db"),
	}
	pullError := collection.Pull()
	db.Collections = append(db.Collections, &collection)
	return &collection, pullError
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
