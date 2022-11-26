package dvbcrud_go

import (
	"fmt"
	"strings"
)

func getSelectFrom(tableName string, idFieldName string) string {
	if idFieldName == "" {
		return fmt.Sprintf("SELECT * FROM %s;", tableName)
	}
	return fmt.Sprintf("SELECT * FROM %s WHERE %s = ?;", tableName, idFieldName)
}

func getInsertInto(tableName string, fields ...string) string {
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);",
		tableName,
		strings.Join(fields, ", "),
		strings.Repeat("?, ", len(fields)-1)+"?")
}

func getUpdate(tableName string, idFieldName string, fields ...string) string {
	for i := range fields {
		fields[i] += " = ?"
	}
	return fmt.Sprintf("UPDATE %s SET (%s) WHERE %s = ?;", tableName, strings.Join(fields, ", "), idFieldName)
}

func getDeleteFrom(tableName string, idFieldName string) string {
	return fmt.Sprintf("DELETE * FROM %s WHERE %s = ?;", tableName, idFieldName)
}
