package main

import (
	"fmt"
	"log"
	"object_db"
	"reflect"
)

type Person struct {
	Age   int
	Name  string
	Title string
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	collection := object_db.Collection{
		FilePath: "example.db",
		DType:    reflect.TypeOf(Person{}),
	}

	john, err := collection.Insert(Person{
		Age:  20,
		Name: "John",
	})
	handleError(err)

	david, err := collection.Insert(Person{
		Age:  30,
		Name: "David",
	})
	handleError(err)

	j := collection.Get(john).(Person)
	fmt.Println(j)

	fmt.Println(collection.Get(david))
}
