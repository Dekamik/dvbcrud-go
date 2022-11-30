package dvbcrud

import (
	"fmt"
)

// paramType separates Column and Value parameter types.
// This is only applicable to prepared statements in Oracle.
type paramType int

const (
	Columns paramType = iota
	Values
)

type sqlParameterGenerator interface {
	GetParamPlaceholders(amount int, typ paramType) ([]string, error)
}

type sqlParameterGeneratorImpl struct {
	sqlParameterGenerator
	dialect SQLDialect
}

// GetParamPlaceholders returns n amount of parameter placeholders as an array of strings.
// The placeholders are formatted according to the chosen dialect.
// (e.g. MySQL-like = ?, PostgreSQL = $1, Oracle = :col1 or :var1)
func (p sqlParameterGeneratorImpl) GetParamPlaceholders(amount int, typ paramType) ([]string, error) {
	placeholders := make([]string, amount)

	switch p.dialect {
	case MySQL, SQLite, ODBC, MariaDB:
		for i := 0; i < amount; i++ {
			placeholders[i] = "?"
		}

	case PostgreSQL:
		for i := 0; i < amount; i++ {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
		}

	case Oracle:
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

func newSQLParamGen(dialect SQLDialect) sqlParameterGenerator {
	return sqlParameterGeneratorImpl{
		dialect: dialect,
	}
}
