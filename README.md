# LightDB
Lightweight object database written in Go

## Insert Example
```go
import (
    "lightdb"
    ...
)

type Person struct {
    Age   int
    Name  string
    Title string
}

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

_john := collection.Get(john).(Person)

```

### Find object in collection using fields map
```go
personSearch := make(map[string]interface{})
personSearch["Name"] = "David"

results := make([]interface{}, 0)
collection.Find(personSearch, &results)

```

### Commit changes to disk
```go
commitError := collection.Commit()
if commitError != nil {
    log.Fatal("Error commiting db changes")
}
```

### Pull changes from disk
```go
func main() {
    collection := lightdb.Collection{
        FilePath: "example.db",
        DType:    reflect.TypeOf(Person{}),
    }
	
    pullError := collection.Pull()
    if pullError != nil {
        log.Fatal("Error pulling db from disk")
    } else {
        log.Println("pulled db from disk, db is up-to-date!")	
    } 
}
```