package lightdb

import (
	"os"
	"path"
	"testing"
)

type StoreTestStruct struct {
	x string
	Y int
	T bool
	Z byte
}

func TestStore(t *testing.T) {
	obj := StoreTestStruct{
		x: "hello",
		Y: 123,
		T: false,
		Z: 0x1,
	}

	test_file_path := path.Join(os.TempDir(), "test_store_write.data")

	writeError := writeObject(test_file_path, obj)
	if writeError != nil {
		t.Fatal(writeError)
	}

	_obj := StoreTestStruct{}
	readError := readObject(test_file_path, &_obj)
	if readError != nil {
		t.Fatal(readError)
	}
}

//func TestStore_JSON(t *testing.T) {
//	test_file_path := path.Join(os.TempDir(), "test_store_json.db")
//	t.Log(test_file_path)
//	db := Collection{Name: "test-collection", FilePath: "test_collection_find.db", DType: reflect.TypeOf(StoreTestStruct{})}
//	db.Insert(StoreTestStruct{Y: 1, T: false})
//	db.Insert(StoreTestStruct{Y: 2, T: true})
//	db.Insert(StoreTestStruct{Y: 3, T: false})
//	db.Insert(StoreTestStruct{Y: 4, T: true})
//	db.Insert(StoreTestStruct{Y: 5, T: false})
//	db.Insert(StoreTestStruct{Y: 6, T: true})
//
//	err := writeJSON(test_file_path, db.CollectionCore.Index)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	_db := Collection{Name: "test-collection", FilePath: "test_collection_find.db", DType: reflect.TypeOf(StoreTestStruct{})}
//
//	err = readJSON(test_file_path, &_db.CollectionCore.Index)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	x := _db.Get(1).(interface{})
//	t.Log(x)
//}
