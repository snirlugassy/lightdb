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

...

func main() {
    collection := lightdb.Collection{
        FilePath: "example.db",
        DType:    reflect.TypeOf(Person{}),
    }
    
    john, err := collection.Insert(Person{
        Age:  20,
        Name: "John",
    })
    handleError(err)
    
    j := collection.Get(john).(Person)
    log.Println(j)
}
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