# dvbcrud-go

Simple generic CRUD wrapper for `sqlx`.

# Usage

```go
package example

import (
    "fmt"
    "github.com/Dekamik/dvbcrud-go"
    "github.com/jmoiron/sqlx"
    "time"
)

type user struct {
    Id        int `db:"UserId"`
    Name      string
    Birthdate time.Time
}

func main() {
    db, _ := sqlx.Connect("mssql", "datasource") // err ignored for brevity
    userRepo := dvbcrud.SqlRepository[user]{
        db:          db,
        tableName:   "Users",
        idFieldName: "UserId",
    }

    users, _ := userRepo.ReadAll() // err ignored for brevity

    for _, u := range users {
        fmt.Printf("Hello, %s!\n", u.Name)
    }
}
```
