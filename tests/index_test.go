package tests

import (
	"github.com/snirlugassy/lightdb"
	"reflect"
	"testing"
)

func TestHashIndex(t *testing.T) {
	collection := lightdb.Collection{FilePath: "test_collection_find.collection", DType: reflect.TypeOf(A{})}
	collection.Insert(A{Name: "a", Age: 1})
	collection.Insert(A{Name: "b", Age: 2})
	collection.Insert(A{Name: "b", Age: 3})
	collection.Insert(A{Name: "c", Age: 4})
	collection.Insert(A{Name: "d", Age: 4})
	collection.Insert(A{Name: "d", Age: 5})

	index := lightdb.HashIndex{Collection: &collection}
	index.Build("Name")
	t.Log(index.Index["d"])
}
