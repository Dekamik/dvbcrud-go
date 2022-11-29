package internal

import (
	"fmt"
	"github.com/dekamik/dvbcrud-go/crudsql"
)

// ParamType separates Column and Value parameter types.
// This is only applicable to prepared statements in Oracle.
type ParamType int

const (
	Columns ParamType = iota
	Values
)

type ParamGen interface {
	GetParamPlaceholders(amount int, typ ParamType) ([]string, error)
}

type paramGenImpl struct {
	ParamGen
	dialect crudsql.SQLDialect
}

// GetParamPlaceholders returns n amount of parameter placeholders as an array of strings.
// The placeholders are formatted according to the chosen dialect.
// (e.g. MySQL-like = ?, PostgreSQL = $1, Oracle = :col1 or :var1)
func (p paramGenImpl) GetParamPlaceholders(amount int, typ ParamType) ([]string, error) {
	placeholders := make([]string, amount)

	switch p.dialect {
	case crudsql.MySQL, crudsql.SQLite, crudsql.ODBC, crudsql.MariaDB:
		for i := 0; i < amount; i++ {
			placeholders[i] = "?"
		}

	case crudsql.PostgreSQL:
		for i := 0; i < amount; i++ {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
		}

	case crudsql.Oracle:
		var name string
		if typ == Columns {
			name = "col"
		} else {
			name = "val"
		}
		if amount == 1 {
			placeholders[0] = ":" + name
			break
		}
		for i := 0; i < amount; i++ {
			placeholders[i] = fmt.Sprintf(":%s%d", name, i+1)
		}

	default:
		return nil, fmt.Errorf("unknown dialect")
	}

	return placeholders, nil
}

func NewSQLParamGen(dialect crudsql.SQLDialect) ParamGen {
	return paramGenImpl{
		dialect: dialect,
	}
}
