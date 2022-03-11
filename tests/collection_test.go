package tests

import (
	"object_db"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type A struct {
	Age  int
	Name string
}

type B struct {
	Name string
}

func TestCollection_Insert(t *testing.T) {
	db_file_path := filepath.Join(os.TempDir(), "test_insert.db")
	db := object_db.Collection{FilePath: db_file_path, DType: reflect.TypeOf(A{})}
	obj := A{Age: 1}
	i, err := db.Insert(&obj)
	if err != nil {
		t.Log(err)
		t.Fatal("Error inserting A")
	}

	_obj := db.Get(i).(*A)
	t.Log(_obj.Age)
	if !reflect.DeepEqual(obj, *_obj) {
		t.Fatal("restored object different from inserted object")
	}
}

func TestCollection_Commit(t *testing.T) {
	db_file_path := filepath.Join(os.TempDir(), "test_commit.db")
	db := object_db.Collection{FilePath: db_file_path, DType: reflect.TypeOf(A{})}
	obj := A{Age: 1}
	_, err := db.Insert(obj)
	if err != nil {
		t.Fatal(err)
	}
	commitError := db.Commit()
	if commitError != nil {
		t.Fatal(commitError)
	}
}

func TestCollection_Pull(t *testing.T) {
	db_file_path := filepath.Join(os.TempDir(), "test_pull.db")
	db := object_db.Collection{FilePath: db_file_path, DType: reflect.TypeOf(A{})}
	obj := &A{Age: 1}
	i, err := db.Insert(obj)
	if err != nil {
		t.Fatal(err)
	}
	commitErr := db.Commit()
	if commitErr != nil {
		t.Fatal(commitErr)
	}

	db2 := object_db.Collection{FilePath: db_file_path, DType: reflect.TypeOf(A{})}
	pullError := db2.Pull()
	if pullError != nil {
		t.Fatal(pullError)
	}
	obj2 := db2.Get(i).(*A)

	if *obj != *obj2 {
		t.Log(obj)
		t.Log(obj2)
		t.Fatal("Failed to commit, pull and restore object")
	}
}

func TestCollection_StrictType(t *testing.T) {
	db := object_db.Collection{FilePath: "test_strict_type.db"}
	obj := A{Age: 1}
	_, errA := db.Insert(obj)
	if errA != nil {
		t.Log(errA)
		t.Fatal("Error inserting A")
	}

	b := B{Name: "hello"}
	_, errB := db.Insert(b)
	if errB == nil {
		t.Fatal(errB)
	}
}
