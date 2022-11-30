package crudsql

type sqlTemplatesMock struct {
	sqlTemplates

	GetSelectReturn    string
	GetSelectAllReturn string
	GetInsertReturn    string
	GetUpdateReturn    string
	GetDeleteReturn    string

	GetInsertError error
	GetUpdateError error
}

func (s sqlTemplatesMock) GetSelect() string {
	return s.GetSelectReturn
}

func (s sqlTemplatesMock) GetSelectAll() string {
	return s.GetSelectAllReturn
}

func (s sqlTemplatesMock) GetInsert(fields []string) (string, error) {
	return s.GetInsertReturn, s.GetInsertError
}

func (s sqlTemplatesMock) GetUpdate(fields []string) (string, error) {
	return s.GetUpdateReturn, s.GetUpdateError
}

func (s sqlTemplatesMock) GetDelete() string {
	return s.GetDeleteReturn
}
