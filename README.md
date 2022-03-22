
# LightDB

[![build-and-test-badge](https://img.shields.io/github/workflow/status/snirlugassy/lightdb/Build%20and%20Test/main)](https://github.com/snirlugassy/lightdb/actions/workflows/build.yml?branch=main)
![license-badge](https://img.shields.io/github/license/snirlugassy/lightdb)
![last-commit-badge](https://img.shields.io/github/last-commit/snirlugassy/lightdb)
[![group-badge](https://img.shields.io/badge/-Group-yellow)](https://groups.google.com/g/lightdb)

Lightweight object database written in Go

For more information join our [Google Group](https://groups.google.com/g/lightdb)

### Installation
`go get github.com/snirlugassy/lightdb`

### Initialize DB instance
```go
import "github.com/snirlugassy/lightdb"

db := lightdb.Database{
    Name: "example-db",
    Path: "/tmp/db",
}
```

### Create collection
```go
collection := db.CreateCollection("users", reflect.TypeOf(User{}))
```

### Insert objects
```go
john, err := collection.Insert(User{ID: 1, Name: "John"})
if err != nil {
    log.Fatal("error inserting john to collection")
} else {
    log.Println("John's ID", john)
}

collection.Insert(User{ID: 2, Name: "David", IsAdmin: true})
collection.Insert(User{ID: 3, Name: "Alfred", IsAdmin: false})
```

### Retrieve objects
```go
user1, found := collection.Get(1).(User)
if !found {
    log.Fatal("John not found")
}
```

### Search objects
```go
personSearch := make(map[string]interface{})
personSearch["Name"] = "David"

results := make([]interface{}, 0)
collection.Find(personSearch, &results)
```

### Update objects
```go
log.Println("Setting John as admin")
user1.IsAdmin = true

log.Println("Updating collection")
collection.Update(user1.ID, user1)
```

### Commit changes (Save to disk)
```go
db.Commit()
```

### Create another DB instance
```go
db2 := lightdb.Database{
    Name: "example-db",
    Path: "/tmp/db",
}

collection2, err := db2.LoadCollection("users", reflect.TypeOf(User{}))
if err != nil {
    log.Fatal(err)
}
collection2.Pull()
_user1, found := collection2.Get(1).(User)
```

### Runtime analysis
#### Insert

![insert-line-chart](https://raw.githubusercontent.com/snirlugassy/lightdb/main/analysis/insert_viz.png)

#### Commit and Pull DB with single collection

![commit-pull-line](https://raw.githubusercontent.com/snirlugassy/lightdb/main/analysis/commit_pull_viz.png)
