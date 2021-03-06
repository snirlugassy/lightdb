package tests

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/snirlugassy/lightdb"
)

type A struct {
	Age  int
	Name string
}

type B struct {
	Title string
}

func TestCollection_Get(t *testing.T) {
	collectionPath := filepath.Join(os.TempDir(), "test_get.collection")
	collection := lightdb.Collection{FilePath: collectionPath, DType: reflect.TypeOf(A{})}
	objID, err := collection.Insert(A{Name: "a", Age: 1})
	if err != nil {
		t.Fatal("failed to insert object")
	}

	obj, found := collection.Get(objID).(A)
	if !found {
		t.Fatal("inserted object not found")
	}

	if obj.Name != "a" || obj.Age != 1 {
		t.Fatal("returned object not matching inserted object")
	}

	obj, found = collection.Get(999).(A)
	if found {
		t.Fatal("not inserted object found")
	}

}

func TestCollection_GetAll(t *testing.T) {
	collectionPath := filepath.Join(os.TempDir(), "test_get_all.collection")
	collection := lightdb.Collection{FilePath: collectionPath, DType: reflect.TypeOf(A{})}
	collection.Insert(A{Name: "a", Age: 1})
	collection.Insert(A{Name: "b", Age: 2})
	collection.Insert(A{Name: "b", Age: 3})
	collection.Insert(A{Name: "c", Age: 4})

	items := collection.GetAll()

	if len(items) != 4 {
		t.Fatal("wrong result size for all items")
	}
}

func TestCollection_Insert(t *testing.T) {
	collectionPath := filepath.Join(os.TempDir(), "test_insert.collection")
	collection := lightdb.Collection{FilePath: collectionPath, DType: reflect.TypeOf(A{})}
	obj := A{Age: 1}
	i, err := collection.Insert(obj)
	if err != nil {
		t.Log(err)
		t.Fatal("Error inserting A")
	}

	_obj, found := collection.Get(i).(A)
	if !found {
		t.Fatal("_obj not found")
	}

	t.Log(_obj.Age)
	if !reflect.DeepEqual(obj, _obj) {
		t.Fatal("restored object different from inserted object")
	}
}

func TestCollection_Commit(t *testing.T) {
	collectionPath := filepath.Join(os.TempDir(), "test_commit.db")
	db := lightdb.Collection{FilePath: collectionPath, DType: reflect.TypeOf(A{})}
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
	collectionPath := filepath.Join(os.TempDir(), "test_pull.db")
	db := lightdb.Collection{FilePath: collectionPath, DType: reflect.TypeOf(A{})}
	obj := A{Age: 1}
	i, err := db.Insert(obj)
	if err != nil {
		t.Fatal(err)
	}
	commitErr := db.Commit()
	if commitErr != nil {
		t.Fatal(commitErr)
	}

	db2 := lightdb.Collection{FilePath: collectionPath, DType: reflect.TypeOf(A{})}
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
	collectionPath := filepath.Join(os.TempDir(), "test_strict_type.collection")
	collection := lightdb.Collection{FilePath: collectionPath, DType: reflect.TypeOf(A{})}
	obj := A{Age: 1}
	_, errA := collection.Insert(obj)
	if errA != nil {
		t.Log(errA)
		t.Fatal("Error inserting A")
	}

	b := B{Title: "hello"}
	_, errB := collection.Insert(b)
	if errB == nil {
		t.Fatal(errB)
	}
}

func TestCollection_Find(t *testing.T) {
	collectionPath := filepath.Join(os.TempDir(), "test_find.collection")
	collection := lightdb.Collection{FilePath: collectionPath, DType: reflect.TypeOf(A{})}
	collection.Insert(A{Name: "a", Age: 1})
	collection.Insert(A{Name: "b", Age: 2})
	collection.Insert(A{Name: "b", Age: 3})
	collection.Insert(A{Name: "c", Age: 4})
	collection.Insert(A{Name: "d", Age: 4})
	collection.Insert(A{Name: "d", Age: 5})

	searchFields := make(map[string]interface{})
	searchFields["Name"] = "b"
	searchFields["Age"] = "20"
	searchFields["Banana"] = 1

	results := make([]interface{}, 0)
	collection.Find(searchFields, &results)
	t.Log(results)

	if len(results) != 2 {
		t.Fatal("wrong results array size")
	}

	searchFields["Name"] = nil
	searchFields["Age"] = 4

	results = make([]interface{}, 0)
	collection.Find(searchFields, &results)
	t.Log(results)

	if len(results) != 2 {
		t.Fatal("wrong results array size")
	}
}

func TestCollection_First(t *testing.T) {
	collectionPath := filepath.Join(os.TempDir(), "test_find_first.collection")
	collection := lightdb.Collection{FilePath: collectionPath, DType: reflect.TypeOf(A{})}
	collection.Insert(A{Name: "a", Age: 1})
	collection.Insert(A{Name: "b", Age: 2})
	collection.Insert(A{Name: "b", Age: 3})
	collection.Insert(A{Name: "c", Age: 4})
	collection.Insert(A{Name: "d", Age: 4})
	collection.Insert(A{Name: "d", Age: 5})

	searchFields := make(map[string]interface{})
	searchFields["Age"] = 5

	var result interface{}
	collection.First(searchFields, &result)

	if result == nil {
		t.Fatal("result is nil")
	}

	if result.(A).Age != 5 {
		t.Fatal("wrong result")
	}
}

func TestCollection_Update(t *testing.T) {
	collectionPath := filepath.Join(os.TempDir(), "test_update.collection")
	collection := lightdb.Collection{FilePath: collectionPath, DType: reflect.TypeOf(A{})}
	i, insertError := collection.Insert(A{Name: "a", Age: 1})
	if insertError != nil {
		t.Fatal(insertError)
	}
	collection.Update(i, A{Name: "b", Age: 2})
	updated, found := collection.Get(i).(A)
	if !found || updated.Name != "b" || updated.Age != 2 {
		t.Fatal("failed to update")
	}
}

func TestCollection_Filter(t *testing.T) {
	collectionPath := filepath.Join(os.TempDir(), "test_filter.collection")
	collection := lightdb.Collection{FilePath: collectionPath, DType: reflect.TypeOf(A{})}
	collection.Insert(A{Name: "a", Age: 1})
	collection.Insert(A{Name: "b", Age: 2})
	collection.Insert(A{Name: "b", Age: 3})
	collection.Insert(A{Name: "c", Age: 4})
	collection.Insert(A{Name: "d", Age: 4})
	collection.Insert(A{Name: "d", Age: 5})

	results := make([]interface{}, 0)
	filter := func(a interface{}) bool { return a.(A).Age > 3 }

	collection.Filter(filter, &results)

	if len(results) != 3 {
		t.Fatal("wrong results array size")
	}
}

func TestCollection_FilterFirst(t *testing.T) {
	collectionPath := filepath.Join(os.TempDir(), "test_filter_first.collection")
	collection := lightdb.Collection{FilePath: collectionPath, DType: reflect.TypeOf(A{})}
	collection.Insert(A{Name: "a", Age: 1})
	collection.Insert(A{Name: "b", Age: 2})
	collection.Insert(A{Name: "b", Age: 3})

	filter := func(a interface{}) bool { return a.(A).Age > 2 }

	var result interface{}
	collection.FilterFirst(filter, &result)

	if result == nil {
		t.Fatal("result is nil")
	}

	if result.(A).Age != 3 || result.(A).Name != "b" {
		t.Fatal("wrong result")
	}
}

func TestCollection_ToJSON(t *testing.T) {
	collectionPath := filepath.Join(os.TempDir(), "test_to_json.collection")
	collection := lightdb.Collection{FilePath: collectionPath, DType: reflect.TypeOf(A{})}
	collection.Insert(A{Name: "a", Age: 1})

	data, err := collection.ToJSON()
	if err != nil {
		t.Fatal(err)
	}

	if len(string(data)) == 0 {
		t.Fatal("empty json returned")
	}
}
