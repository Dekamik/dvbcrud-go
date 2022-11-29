package dvbcrud

import (
	"fmt"
	"strings"
)

type SQLDialect int

const (
	MySQL SQLDialect = iota
	PostgreSQL
	Oracle
	SQLite
	ODBC
	MariaDB
)

type paramType int

const (
	Columns paramType = iota
	Values
)

func getParamPlaceholders(dialect SQLDialect, typ paramType, amount int) ([]string, error) {
	placeholders := make([]string, amount)

	switch dialect {
	case MySQL, SQLite, ODBC, MariaDB:
		for i := 0; i < amount; i++ {
			placeholders[i] = "?"
		}
		break

	case PostgreSQL:
		for i := 0; i < amount; i++ {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
		}
		break

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
		break

	default:
		return nil, fmt.Errorf("unknown dialect")
	}

	return placeholders, nil
}

func getSelectFromStmt(dialect SQLDialect, tableName string, idFieldName string, fields ...string) (string, error) {
	placeholders, err := getParamPlaceholders(dialect, Columns, 1)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("SELECT %s FROM %s WHERE %s = %s;",
		strings.Join(fields, ", "),
		tableName,
		idFieldName,
		placeholders[0]), nil
}

func getSelectAllStmt(tableName string, fields ...string) string {
	return fmt.Sprintf("SELECT %s FROM %s;", strings.Join(fields, ", "), tableName)
}

func getInsertIntoStmt(dialect SQLDialect, tableName string, fields ...string) (string, error) {
	placeholders, err := getParamPlaceholders(dialect, Values, len(fields))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);",
		tableName,
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", ")), nil
}

func getUpdateStmt(dialect SQLDialect, tableName string, idFieldName string, fields ...string) (string, error) {
	columnPlaceholders, err := getParamPlaceholders(dialect, Columns, 1)
	if err != nil {
		return "", err
	}

	valuePlaceholders, err := getParamPlaceholders(dialect, Values, len(fields))
	if err != nil {
		return "", err
	}

	for i := range fields {
		fields[i] += " = " + valuePlaceholders[i]
	}

	return fmt.Sprintf("UPDATE %s SET (%s) WHERE %s = %s;",
		tableName,
		strings.Join(fields, ", "),
		idFieldName,
		columnPlaceholders[0]), nil
}

func getDeleteFromStmt(dialect SQLDialect, tableName string, idFieldName string) (string, error) {
	placeholder, err := getParamPlaceholders(dialect, Columns, 1)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("DELETE FROM %s WHERE %s = %s;",
		tableName,
		idFieldName,
		placeholder[0]), nil
}
