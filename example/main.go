package main

import (
	"log"
	lightdb "object_db"
	"reflect"
)

type Person struct {
	Age   int
	Name  string
	Title string
}

func main() {
	db := lightdb.Database{
		Name: "example-db",
		Path: "/tmp/db",
	}

	collection := db.CreateCollection("person", reflect.TypeOf(Person{}))

	john, err := collection.Insert(Person{Age: 20, Name: "John"})
	if err != nil {
		log.Fatal("error inserting john to collection")
	}

	collection.Insert(Person{Age: 30, Name: "David"})

	j := collection.Get(john).(Person)
	log.Println(j)

	personSearch := make(map[string]interface{})
	personSearch["Name"] = "David"

	results := make([]interface{}, 0)
	collection.Find(personSearch, &results)
	log.Println(results)
}
