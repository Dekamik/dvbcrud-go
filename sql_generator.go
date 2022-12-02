package dvbcrud

import (
    "fmt"
    "strings"
)

type sqlGenerator interface {
    // GenerateSelect generates and returns a SELECT statement (WHERE ID)
    GenerateSelect(table string, idField string, fields []string) (string, error)

    // GenerateSelectAll generates and returns a SELECT statement (all rows)
    GenerateSelectAll(table string, fields []string) string

    // GenerateInsert generates and returns an INSERT INTO statement
    GenerateInsert(table string, fields []string) (string, error)

    // GenerateUpdate generates and returns an UPDATE statement
    GenerateUpdate(table string, idField string, fields []string) (string, error)

    // GenerateDelete returns DELETE FROM <table> WHERE <id> = ?
    GenerateDelete(table string, idField string) (string, error)
}

type sqlGeneratorImpl struct {
    sqlGenerator
    sqlParamGen sqlParameterGenerator
}

func (s sqlGeneratorImpl) GenerateSelect(table string, idField string, fields []string) (string, error) {
    placeholders, err := s.sqlParamGen.GetParamPlaceholders(1, Columns)
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("SELECT %s FROM %s WHERE %s = %s",
        strings.Join(fields, ", "),
        table,
        idField,
        placeholders[0]), nil
}

func (s sqlGeneratorImpl) GenerateSelectAll(table string, fields []string) string {
    return fmt.Sprintf("SELECT %s FROM %s", strings.Join(fields, ", "), table)
}

func (s sqlGeneratorImpl) GenerateInsert(table string, fields []string) (string, error) {
    placeholders, err := s.sqlParamGen.GetParamPlaceholders(len(fields), Values)
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
        table,
        strings.Join(fields, ", "),
        strings.Join(placeholders, ", ")), nil
}

func (s sqlGeneratorImpl) GenerateUpdate(table string, idField string, fields []string) (string, error) {
    columnPlaceholders, err := s.sqlParamGen.GetParamPlaceholders(1, Columns)
    if err != nil {
        return "", err
    }

    valuePlaceholders, err := s.sqlParamGen.GetParamPlaceholders(len(fields), Values)
    if err != nil {
        return "", err
    }

    f := fields
    for i := range f {
        f[i] += " = " + valuePlaceholders[i]
    }

    return fmt.Sprintf("UPDATE %s SET (%s) WHERE %s = %s",
        table,
        strings.Join(f, ", "),
        idField,
        columnPlaceholders[0]), nil
}

func (s sqlGeneratorImpl) GenerateDelete(table string, idField string) (string, error) {
    placeholder, err := s.sqlParamGen.GetParamPlaceholders(1, Columns)
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("DELETE FROM %s WHERE %s = %s",
        table,
        idField,
        placeholder[0]), nil
}

func newSQLGenerator(paramGen sqlParameterGenerator) sqlGenerator {
    return &sqlGeneratorImpl{
        sqlParamGen: paramGen,
    }
}
