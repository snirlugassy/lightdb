package lightdb

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
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

func (db *Database) DirName() string {
	path, _ := os.Getwd()
	pathSlice := regexp.MustCompile("[^0-9A-Za-z_]").Split(path, -1)
	currentDir := pathSlice[len(pathSlice)-1]
	return currentDir

}

func (db *Database) CalculateFolderSize() (dirsize int64, err error) {

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

func (db *Database) Info() (info map[string]interface{}) {
	folder := db.DirName()

	fileSize, _ := db.CalculateFolderSize()

	info = make(map[string]interface{})

	info["Folder"] = folder
	info["File Size(bytes)"] = fileSize
	info["Database name"] = db.Name
	return

}
