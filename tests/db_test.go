package tests

import (
	"os"
	"reflect"
	"testing"

	"github.com/snirlugassy/lightdb"
)

type Player struct {
	Name  string
	Level int
}

func TestDatabase_Init(t *testing.T) {
	db := lightdb.Database{Name: "testdb"}
	initError := db.Init()
	if !os.IsNotExist(initError) {
		t.Fatal("db was initialized without path")
	}

	db = lightdb.Database{Name: "testdb", Path: "/tmp/lightdb"}
	initError = db.Init()
	if initError != nil {
		t.Fatal(initError)
	}
}

func TestDatabase_InitCollection(t *testing.T) {
	db := lightdb.Database{Name: "testdb"}
	db.Init()
	db.InitCollection("players-collection", reflect.TypeOf(Player{}))

	if db.Collections == nil || len(db.Collections) != 1 {
		t.Fatal("db create collection failed failed")
	}
}
