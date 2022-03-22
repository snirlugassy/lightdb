package tests

import (
	"github.com/snirlugassy/lightdb"
	"reflect"
	"testing"
)

type Player struct {
	Name  string
	Level int
}

func TestDatabase_Init(t *testing.T) {
	db := lightdb.Database{Name: "testdb"}
	db.Init()
	if db.Collections == nil || db.Name != "testdb" {
		t.Fatal("db init failed")
	}
}

func TestDatabase_CreateCollection(t *testing.T) {
	db := lightdb.Database{Name: "testdb"}
	db.Init()
	db.CreateCollection("players-collection", reflect.TypeOf(Player{}))

	if db.Collections == nil || len(db.Collections) != 1 {
		t.Fatal("db create collection failed failed")
	}
}
