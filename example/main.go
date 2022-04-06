package main

import (
	"fmt"
	"log"
	"reflect"

	"github.com/snirlugassy/lightdb"
)

type User struct {
	ID      int
	Name    string
	IsAdmin bool
}

func main() {
	db := lightdb.Database{
		Name: "example-db",
		Path: "/tmp/db",
	}

	collection := db.InitCollection("users", reflect.TypeOf(User{}))

	john, err := collection.Insert(User{ID: 1, Name: "John"})
	if err != nil {
		log.Fatal("error inserting john to collection")
	} else {
		log.Println("John's ID", john)
	}

	collection.Insert(User{ID: 2, Name: "David", IsAdmin: true})
	collection.Insert(User{ID: 3, Name: "Alfred", IsAdmin: false})

	user1, found := collection.Get(1).(User)
	if !found {
		log.Fatal("John not found")
	} else {
		log.Println("user1 Name before update: ", user1.Name)
		log.Println("user1 IsAdmin before update: ", user1.IsAdmin)
	}

	log.Println("Setting John as admin")
	user1.IsAdmin = true

	log.Println("Updating collection")
	collection.Update(user1.ID, user1)

	user1, found = collection.Get(1).(User)
	if !found {
		log.Fatal("John not found after update")
	} else {
		log.Println("user1 Name after update: ", user1.Name)
		log.Println("user1 IsAdmin after update: ", user1.IsAdmin)
	}

	log.Println("Saving DB Changes")
	db.Commit()

	log.Println("Creating second DB instance")
	db2 := lightdb.Database{
		Name: "example-db",
		Path: "/tmp/db",
	}

	log.Println("Loading users collection")
	collection2, err := db2.LoadCollection("users", reflect.TypeOf(User{}))
	if err != nil {
		log.Fatal(err)
	}
	collection2.Pull()

	_user1, found := collection2.Get(1).(User)
	log.Println("Found user:", found)
	log.Println(_user1)

	personSearch := make(map[string]interface{})
	personSearch["Name"] = "David"

	results := make([]interface{}, 0)
	collection.Find(personSearch, &results)
	log.Println(results)

	mydir := collection.DirName()
	fmt.Println(mydir)

	myFolder, err := collection.CalculateFolderSize()
	log.Println(myFolder, err)

	myInfo := collection.Info()
	log.Println(myInfo)

	myDbInfo := db.Info()
	log.Println(myDbInfo)
}
