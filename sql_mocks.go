package dvbcrud

import "reflect"

type paramGenMock struct {
	sqlParameterGenerator
	GetParamPlaceholdersMock func(amount int, typ paramType) ([]string, error)
}

func (p paramGenMock) GetParamPlaceholders(amount int, typ paramType) ([]string, error) {
	return p.GetParamPlaceholdersMock(amount, typ)
}

type sqlGeneratorMock struct {
	sqlGenerator

	generateSelectMock    func(table string, idField string, fields []string) (string, error)
	generateSelectAllMock func(table string, fields []string) string
	generateInsertMock    func(table string, fields []string) (string, error)
	generateUpdateMock    func(table string, idField string, fields []string) (string, error)
	generateDeleteMock    func(table string, idField string) (string, error)
}

func (s sqlGeneratorMock) generateSelect(table string, idField string, fields []string) (string, error) {
	return s.generateSelectMock(table, idField, fields)
}

func (s sqlGeneratorMock) generateSelectAll(table string, fields []string) string {
	return s.generateSelectAllMock(table, fields)
}

func (s sqlGeneratorMock) generateInsert(table string, fields []string) (string, error) {
	return s.generateInsertMock(table, fields)
}

func (s sqlGeneratorMock) generateUpdate(table string, idField string, fields []string) (string, error) {
	return s.generateUpdateMock(table, idField, fields)
}

func (s sqlGeneratorMock) generateDelete(table string, idField string) (string, error) {
	return s.generateDeleteMock(table, idField)
}

type sqlTemplatesMock struct {
	sqlTemplates

	GetSelectMock    func() string
	GetSelectAllMock func() string
	GetInsertMock    func(fields []string) (string, error)
	GetUpdateMock    func(fields []string) (string, error)
	GetDeleteMock    func() string
}

func (s sqlTemplatesMock) GetSelect() string {
	return s.GetSelectMock()
}

func (s sqlTemplatesMock) GetSelectAll() string {
	return s.GetSelectAllMock()
}

func (s sqlTemplatesMock) GetInsert(fields []string) (string, error) {
	return s.GetInsertMock(fields)
}

func (s sqlTemplatesMock) GetUpdate(fields []string) (string, error) {
	return s.GetUpdateMock(fields)
}

func (s sqlTemplatesMock) GetDelete() string {
	return s.GetDeleteMock()
}

type structParserMock struct {
	StructParser

	ParseFieldNamesMock func(typ reflect.Type) ([]string, error)
	ParsePropertiesMock func(model any, idFieldName string) ([]string, []any, error)
}

func (s structParserMock) ParseFieldNames(typ reflect.Type) ([]string, error) {
	return s.ParseFieldNamesMock(typ)
}

func (s structParserMock) ParseProperties(model any, idFieldName string) ([]string, []any, error) {
	return s.ParsePropertiesMock(model, idFieldName)
}
