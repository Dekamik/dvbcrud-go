package dvbcrud_go

type Repository[TModel any] interface {
	// SqlCreate inserts a new row in the database with the data in model.
	SqlCreate(model TModel) error

	// SqlRead finds and returns the row with the corresponding id as a *TModel.
	SqlRead(id any) (*TModel, error)

	// SqlReadAll finds and returns all rows for this table in the database as a []TModel.
	SqlReadAll() ([]TModel, error)

	// SqlUpdate finds the row with the corresponding id and updates its
	// columns according to the data in model.
	SqlUpdate(id any, model TModel) error

	// SqlDelete deletes the row with the corresponding id.
	SqlDelete(id any) error
}
