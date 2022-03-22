package tests

import (
	"github.com/snirlugassy/lightdb"
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
	Title string
}

func TestCollection_Insert(t *testing.T) {
	db_file_path := filepath.Join(os.TempDir(), "test_insert.db")
	db := lightdb.Collection{FilePath: db_file_path, DType: reflect.TypeOf(A{})}
	obj := A{Age: 1}
	i, err := db.Insert(obj)
	if err != nil {
		t.Log(err)
		t.Fatal("Error inserting A")
	}

	_obj, found := db.Get(i).(A)
	if !found {
		t.Fatal("_obj not found")
	}

	t.Log(_obj.Age)
	if !reflect.DeepEqual(obj, _obj) {
		t.Fatal("restored object different from inserted object")
	}
}

func TestCollection_Commit(t *testing.T) {
	db_file_path := filepath.Join(os.TempDir(), "test_commit.db")
	db := lightdb.Collection{FilePath: db_file_path, DType: reflect.TypeOf(A{})}
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
	db := lightdb.Collection{FilePath: db_file_path, DType: reflect.TypeOf(A{})}
	obj := A{Age: 1}
	i, err := db.Insert(obj)
	if err != nil {
		t.Fatal(err)
	}
	commitErr := db.Commit()
	if commitErr != nil {
		t.Fatal(commitErr)
	}

	db2 := lightdb.Collection{FilePath: db_file_path, DType: reflect.TypeOf(A{})}
	pullError := db2.Pull()
	if pullError != nil {
		t.Fatal(pullError)
	}
	obj2 := db2.Get(i).(A)

	if !reflect.DeepEqual(obj, obj2) {
		t.Log(obj)
		t.Log(obj2)
		t.Fatal("Failed to commit, pull and restore object")
	}
}

func TestCollection_StrictType(t *testing.T) {
	db := lightdb.Collection{FilePath: "test_strict_type.db", DType: reflect.TypeOf(A{})}
	obj := A{Age: 1}
	_, errA := db.Insert(obj)
	if errA != nil {
		t.Log(errA)
		t.Fatal("Error inserting A")
	}

	b := B{Title: "hello"}
	_, errB := db.Insert(b)
	if errB == nil {
		t.Fatal(errB)
	}
}

func TestCollection_Find(t *testing.T) {
	db := lightdb.Collection{FilePath: "test_collection_find.db", DType: reflect.TypeOf(A{})}
	db.Insert(A{Name: "a", Age: 1})
	db.Insert(A{Name: "b", Age: 2})
	db.Insert(A{Name: "b", Age: 3})
	db.Insert(A{Name: "c", Age: 4})
	db.Insert(A{Name: "d", Age: 4})
	db.Insert(A{Name: "d", Age: 5})

	searchFields := make(map[string]interface{})
	searchFields["Name"] = "b"
	searchFields["Age"] = "20"
	searchFields["Banana"] = 1

	results := make([]interface{}, 0)
	db.Find(searchFields, &results)
	t.Log(results)

	if len(results) != 2 {
		t.Fatal("wrong results array size")
	}

	searchFields["Name"] = nil
	searchFields["Age"] = 4

	results = make([]interface{}, 0)
	db.Find(searchFields, &results)
	t.Log(results)

	if len(results) != 2 {
		t.Fatal("wrong results array size")
	}
}

func TestCollection_First(t *testing.T) {
	db := lightdb.Collection{FilePath: "test_collection_find.db", DType: reflect.TypeOf(A{})}
	db.Insert(A{Name: "a", Age: 1})
	db.Insert(A{Name: "b", Age: 2})
	db.Insert(A{Name: "b", Age: 3})
	db.Insert(A{Name: "c", Age: 4})
	db.Insert(A{Name: "d", Age: 4})
	db.Insert(A{Name: "d", Age: 5})

	searchFields := make(map[string]interface{})
	searchFields["Age"] = 5

	var result interface{}
	db.First(searchFields, &result)
	t.Log(result)
	if result == nil {
		t.Fatal("result is nil")
	}
	if result.(A).Age != 5 {
		t.Fatal("wrong result")
	}
}

func TestCollection_Update(t *testing.T) {
	db := lightdb.Collection{FilePath: "test_collection_find.db", DType: reflect.TypeOf(A{})}
	i, insertError := db.Insert(A{Name: "a", Age: 1})
	if insertError != nil {
		t.Fatal(insertError)
	}
	db.Update(i, A{Name: "b", Age: 2})
	updated, found := db.Get(i).(A)
	if !found || updated.Name != "b" || updated.Age != 2 {
		t.Fatal("failed to update")
	}
}
