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

	for id, obj := range index.Collection.Index {
		fieldValue := reflect.ValueOf(obj).FieldByName(field).Interface()
		index.Index[fieldValue] = append(index.Index[fieldValue], id)
	}
}

//Collection
//	1 -> p1
//	2 -> p2
//	3 -> p3

//Index
//	x -> i1, i2
//	y -> i3, i4
