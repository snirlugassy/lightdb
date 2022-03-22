package tests

import (
	"github.com/snirlugassy/lightdb"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestHashIndex(t *testing.T) {
	collectionPath := filepath.Join(os.TempDir(), "test_hash_index.collection")
	collection := lightdb.Collection{FilePath: collectionPath, DType: reflect.TypeOf(A{})}
	collection.Insert(A{Name: "a", Age: 1})
	collection.Insert(A{Name: "b", Age: 2})
	collection.Insert(A{Name: "b", Age: 3})
	collection.Insert(A{Name: "c", Age: 4})
	collection.Insert(A{Name: "d", Age: 4})
	collection.Insert(A{Name: "d", Age: 5})

	// NAME INDEX TEST
	nameIndex := lightdb.HashIndex{}
	nameIndexBuildError := nameIndex.Build(&collection, "Name")
	if nameIndexBuildError != nil {
		t.Fatal(nameIndexBuildError)
	}

	result := nameIndex.Get("d")
	if result == nil || len(result) != 2 {
		t.Fatal("expected != result in hash nameIndex")
	}

	if nameIndex.Get("blahblah") != nil {
		t.Fatal("nameIndex returns non-nil for non-existing key")
	}

	// AGE INDEX TEST
	ageIndex := lightdb.HashIndex{}
	ageIndexBuildError := ageIndex.Build(&collection, "Age")
	if ageIndexBuildError != nil {
		t.Fatal(ageIndexBuildError)
	}

	result = ageIndex.Get(4)
	if result == nil || len(result) != 2 {
		t.Log(result)
		t.Fatal("expected != result in hash ageIndex")
	}

	if ageIndex.Get("blahblah") != nil {
		t.Fatal("nameIndex returns non-nil for non-existing key")
	}
}
