package dvbcrud

import (
	"fmt"
	"reflect"
)

func parseFieldNames(typ reflect.Type) ([]string, error) {
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type must be a kind of struct")
	}

	numFields := typ.NumField()
	fields := make([]string, numFields)

	for i := 0; i < numFields; i++ {
		name := typ.Field(i).Tag.Get("db")
		if name == "" {
			return nil, fmt.Errorf("%s.%s lacks a db tag", typ.Name(), typ.Field(i).Name)
		}

		fields[i] = name
	}

	return fields, nil
}

// parseProperties reads the struct type T and returns its fields
// and values as two slices. The slices are synchronized which means each
// field and its corresponding value share the same index in both slices.
// Specifying idFieldName filters out that field, which is necessary in
// INSERTS and UPDATES.
func parseProperties(model any, idFieldName string) ([]string, []any, error) {
	val := reflect.ValueOf(model)
	if val.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("model must be a struct type")
	}

	numField := val.NumField()
	fields := make([]string, numField-1)
	values := make([]any, numField-1)

	index := 0
	for i := 0; i < numField; i++ {
		name := val.Type().Field(i).Tag.Get("db")
		if name == "" {
			return nil, nil, fmt.Errorf("%s.%s lacks a db tag", val.Type().Name(), val.Type().Field(i).Name)
		}
		if name == idFieldName {
			continue
		}

		fields[index] = name
		values[index] = val.Field(i).Interface()
		index++
	}

	return fields, values, nil
}
