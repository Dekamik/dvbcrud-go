package dvbcrud

import (
	"fmt"
	"strings"
)

func getSelectFromStmt(tableName string, idFieldName string, fields ...string) string {
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s = ?;", strings.Join(fields, ", "), tableName, idFieldName)
}

func getSelectAllStmt(tableName string, fields ...string) string {
	return fmt.Sprintf("SELECT %s FROM %s;", strings.Join(fields, ", "), tableName)
}

func getInsertIntoStmt(tableName string, fields ...string) string {
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);",
		tableName,
		strings.Join(fields, ", "),
		strings.Repeat("?, ", len(fields)-1)+"?")
}

func getUpdateStmt(tableName string, idFieldName string, fields ...string) string {
	for i := range fields {
		fields[i] += " = ?"
	}
	return fmt.Sprintf("UPDATE %s SET (%s) WHERE %s = ?;", tableName, strings.Join(fields, ", "), idFieldName)
}

func getDeleteFromStmt(tableName string, idFieldName string) string {
	return fmt.Sprintf("DELETE FROM %s WHERE %s = ?;", tableName, idFieldName)
}
