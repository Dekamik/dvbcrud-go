package crudsql

type paramGenMock struct {
	paramGen
	GetParamPlaceholdersMock func(amount int, typ paramType) ([]string, error)
}

func (p paramGenMock) GetParamPlaceholders(amount int, typ paramType) ([]string, error) {
	return p.GetParamPlaceholdersMock(amount, typ)
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
