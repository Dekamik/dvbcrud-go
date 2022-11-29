package crudsql

import (
	"fmt"
	"strings"
)

// SQLDialect denotes the different dialects which define placeholders differently.
type SQLDialect int

const (
	MySQL SQLDialect = iota
	PostgreSQL
	Oracle
	SQLite
	ODBC
	MariaDB
)

// paramType separates Column and Value parameter types.
// This is only applicable to prepared statements in Oracle.
type paramType int

const (
	Columns paramType = iota
	Values
)

// getParamPlaceholders returns n amount of parameter placeholders as an array of strings.
// The placeholders are formatted according to the chosen dialect.
// (e.g. MySQL-like = ?, PostgreSQL = $1, Oracle = :col1 or :var1)
func getParamPlaceholders(dialect SQLDialect, typ paramType, amount int) ([]string, error) {
	placeholders := make([]string, amount)

	switch dialect {
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

// getSelectFromStmt returns SELECT <column>... FROM <table> WHERE <id> = ?
func getSelectFromStmt(dialect SQLDialect, tableName string, idFieldName string, fields ...string) (string, error) {
	placeholders, err := getParamPlaceholders(dialect, Columns, 1)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("SELECT %s FROM %s WHERE %s = %s",
		strings.Join(fields, ", "),
		tableName,
		idFieldName,
		placeholders[0]), nil
}

// getSelectAllStmt returns SELECT <column>... FROM <table>
func getSelectAllStmt(tableName string, fields ...string) string {
	return fmt.Sprintf("SELECT %s FROM %s", strings.Join(fields, ", "), tableName)
}

// getInsertIntoStmt returns INSERT INTO <table> (<column>...) VALUES (?...)
func getInsertIntoStmt(dialect SQLDialect, tableName string, fields ...string) (string, error) {
	placeholders, err := getParamPlaceholders(dialect, Values, len(fields))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", ")), nil
}

// getUpdateStmt returns UPDATE <table> SET <column> = ?... WHERE <id> = ?
func getUpdateStmt(dialect SQLDialect, tableName string, idFieldName string, fields ...string) (string, error) {
	columnPlaceholders, err := getParamPlaceholders(dialect, Columns, 1)
	if err != nil {
		return "", err
	}

	// Not handling this error because the same code is run above
	valuePlaceholders, _ := getParamPlaceholders(dialect, Values, len(fields))

	for i := range fields {
		fields[i] += " = " + valuePlaceholders[i]
	}

	return fmt.Sprintf("UPDATE %s SET (%s) WHERE %s = %s",
		tableName,
		strings.Join(fields, ", "),
		idFieldName,
		columnPlaceholders[0]), nil
}

// getDeleteFromStmt returns DELETE FROM <table> WHERE <id> = ?
func getDeleteFromStmt(dialect SQLDialect, tableName string, idFieldName string) (string, error) {
	placeholder, err := getParamPlaceholders(dialect, Columns, 1)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("DELETE FROM %s WHERE %s = %s",
		tableName,
		idFieldName,
		placeholder[0]), nil
}
