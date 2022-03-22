package lightdb

import (
	"log"
	"reflect"
)

type HashIndex struct {
	Index      map[interface{}][]int
	Collection *Collection
	field      string
}

func (index *HashIndex) Build(field string) {
	_, found := index.Collection.DType.FieldByName(field)
	if !found {
		log.Fatal("field " + field + " does not exists")
	} else {
		log.Println("field found!")
	}

	index.Index = make(map[interface{}][]int)

	for id, obj := range index.Collection.Index {
		fieldValue := reflect.ValueOf(obj).FieldByName(field).Interface()
		index.Index[fieldValue] = append(index.Index[fieldValue], id)
	}
}
