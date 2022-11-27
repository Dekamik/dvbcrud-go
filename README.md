# dvbcrud-go

Simple and generic wrapper for `github.com/jmoiron/sqlx` that handle common CRUD queries as well as conversion between 
rows and struct types. It relies on formatting and reflection to achieve this.

This module is useful for rapidly developing 

# Usage

Below is a simple example on how to use `dvbcrud`. 

```go
package example

import (
    "fmt"
    "github.com/Dekamik/dvbcrud-go"
    "github.com/jmoiron/sqlx"
    "time"
)

type User struct {
    Id        uint64    `db:"UserId"`
    Name      string    `db:"Name"`
    Birthdate time.Time `db:"Birthdate"`
}

func main() {
    // All errors are ignored for brevity
    db, _ := sqlx.Connect("mssql", "datasource")
    _ = db.Ping()
    userRepo, _ := dvbcrud.NewSql[User](db, "Users", "UserId")

    users, _ := userRepo.ReadAll()

    for _, u := range users {
        fmt.Printf("Hello, %s!\n", u.Name)
    }
}
```
