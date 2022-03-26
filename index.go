package lightdb

import (
	"errors"
	"reflect"
)

type DBIndex struct {
	_index map[interface{}][]int
}

type HashIndex struct {
	DBIndex
}

func (index *HashIndex) Build(collection *Collection, field string) error {
	_, found := collection.DType.FieldByName(field)
	if !found {
		return errors.New("field " + field + " does not exists")
	}

	index._index = make(map[interface{}][]int)

	for id, obj := range collection.Index {
		fieldValue := reflect.ValueOf(obj).FieldByName(field).Interface()
		index._index[fieldValue] = append(index._index[fieldValue], id)
	}

	return nil
}

func (index *HashIndex) Get(v interface{}) []int {
	return index._index[v]
}

type BTreeIndex struct {
	DBIndex
}

func (index *BTreeIndex) Build(collection *Collection, field string) error { return nil }
func (index *BTreeIndex) Get(v interface{}) []int                          { return nil }
