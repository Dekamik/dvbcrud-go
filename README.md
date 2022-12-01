# dvbcrud-go

Simple and generic wrapper for `github.com/jmoiron/sqlx` that handle common CRUD queries as well as conversion between 
rows and struct types. It relies on formatting and reflection to achieve this.

This module is useful for rapidly integrating databases in your application, especially as part of a microservice
architecture.

# Usage

Below is a simple example on how to use `dvbcrud`.

```go
package example

import (
    "fmt"
    "github.com/dekamik/dvbcrud-go"
    "github.com/jmoiron/sqlx"
    "time"
)

type User struct {
    ID        uint64    `db:"user_id"`
    Name      string    `db:"name"`
    Birthdate time.Time `db:"birthdate"`
}

func main() {
    // All errors are ignored for brevity
    DB, _ := sqlx.Connect("mssql", "datasource")
    _ = DB.Ping()
    config := dvbcrud.SQLRepositoryConfig{
        dialect: dvbcrud.MySQL,
        table: "Users",
        idField: "UserId",
    }
    userRepo, _ := dvbcrud.NewSQLRepository[User](DB, config)

    seed := []User{
        {
            Name:      "Winston",
            Birthdate: time.Date(1984, time.April, 3, 12, 59, 3, 0, time.UTC),
        },
        {
            Name:      "Julia",
            Birthdate: time.Date(1985, time.April, 3, 12, 59, 3, 0, time.UTC),
        },
    }

    for _, user := range seed {
        _ = userRepo.Create(user)
    }

    results, _ := userRepo.ReadAll()

    for _, result := range results {
        fmt.Printf("Hello, %s!\n", result.Name)
    }
}
```
