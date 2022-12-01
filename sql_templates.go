package dvbcrud

type sqlTemplates interface {
	// GetSelect returns the SELECT statement (WHERE ID)
	GetSelect() string

	// GetSelectAll returns the SELECT statement (all rows)
	GetSelectAll() string

	// GetInsert generates and returns an INSERT INTO statement
	GetInsert(fields []string) (string, error)

	// GetUpdate generates and returns an UPDATE statement
	GetUpdate(fields []string) (string, error)

	// GetDelete returns the DELETE statement
	GetDelete() string
}

type sqlTemplatesImpl struct {
	sqlTemplates
	sqlGen    sqlGenerator
	tableName string
	idField   string

	selectSql    string
	selectAllSql string
	deleteSql    string
}

func (s sqlTemplatesImpl) GetSelect() string {
	return s.selectSql
}

func (s sqlTemplatesImpl) GetSelectAll() string {
	return s.selectAllSql
}

func (s sqlTemplatesImpl) GetInsert(fields []string) (string, error) {
	return s.sqlGen.GenerateInsert(s.tableName, fields)
}

func (s sqlTemplatesImpl) GetUpdate(fields []string) (string, error) {
	return s.sqlGen.GenerateUpdate(s.tableName, s.idField, fields)
}

func (s sqlTemplatesImpl) GetDelete() string {
	return s.deleteSql
}

// newSQLTemplates pre-generates the SELECT, SELECT ALL and DELETE statement and returns a struct containing the templates.
func newSQLTemplates(sqlGen sqlGenerator, tableName string, idField string, allFields []string) (sqlTemplates, error) {
	selectSql, err := sqlGen.GenerateSelect(tableName, idField, allFields)
	if err != nil {
		return nil, err
	}

	selectAllSql := sqlGen.GenerateSelectAll(tableName, allFields)

	deleteSql, err := sqlGen.GenerateDelete(tableName, idField)
	if err != nil {
		return nil, err
	}

	sqlTemp := sqlTemplatesImpl{
		sqlGen:    sqlGen,
		tableName: tableName,
		idField:   idField,
	}

	sqlTemp.selectSql = selectSql
	sqlTemp.selectAllSql = selectAllSql
	sqlTemp.deleteSql = deleteSql

	return &sqlTemp, nil
}
