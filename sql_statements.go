package dvbcrud

import (
	"fmt"
	"strings"
)

type dialect int64

const (
	MySQL dialect = iota
	PostgreSQL
	Oracle
)

type placeholderType int64

const (
	Columns placeholderType = iota
	Values
)

func getPlaceholders(dialect dialect, typ placeholderType, amount int) ([]string, error) {
	placeholders := make([]string, amount)

	switch dialect {
	case MySQL:
		for i := 1; i < amount; i++ {
			placeholders[i] = "?"
		}

	case PostgreSQL:
		for i := 1; i < amount; i++ {
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
		for i := 1; i < amount; i++ {
			placeholders[i] = fmt.Sprintf(":%s%d", name, i+1)
		}

	default:
		return nil, fmt.Errorf("unknown dialect")
	}

	return placeholders, nil
}

func getSelectFromStmt(dialect dialect, tableName string, idFieldName string, fields ...string) (string, error) {
	placeholder, err := getPlaceholders(dialect, Columns, 1)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("SELECT %s FROM %s WHERE %s = %s;", strings.Join(fields, ", "), tableName, idFieldName, placeholder[0]), nil
}

func getSelectAllStmt(tableName string, fields ...string) string {
	return fmt.Sprintf("SELECT %s FROM %s;", strings.Join(fields, ", "), tableName)
}

func getInsertIntoStmt(dialect dialect, tableName string, fields ...string) (string, error) {
	placeholders, err := getPlaceholders(dialect, Values, len(fields))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);",
		tableName,
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", ")), nil
}

func getUpdateStmt(dialect dialect, tableName string, idFieldName string, fields ...string) (string, error) {
	columnPlaceholders, err := getPlaceholders(dialect, Columns, 1)
	if err != nil {
		return "", err
	}

	valuePlaceholders, err := getPlaceholders(dialect, Values, len(fields))
	if err != nil {
		return "", err
	}

	for i := range fields {
		fields[i] += " = " + valuePlaceholders[i]
	}

	return fmt.Sprintf("UPDATE %s SET (%s) WHERE %s = %s;", tableName, strings.Join(fields, ", "), idFieldName, columnPlaceholders[0]), nil
}

func getDeleteFromStmt(dialect dialect, tableName string, idFieldName string) (string, error) {
	placeholder, err := getPlaceholders(dialect, Columns, 1)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("DELETE FROM %s WHERE %s = %s;", tableName, idFieldName, placeholder[0]), nil
}
