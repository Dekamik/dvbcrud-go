package dvbcrud_go

type Repository[TModel any] interface {
	// Create inserts a new row in the database with the data in model.
	Create(model TModel) error

	// Read finds and returns the row with the corresponding id as a *TModel.
	Read(id any) (*TModel, error)

	// ReadAll finds and returns all rows for this table in the database as a []TModel.
	ReadAll() ([]TModel, error)

	// Update finds the row with the corresponding id and updates its
	// columns according to the data in model.
	Update(id any, model TModel) error

	// Delete deletes the row with the corresponding id.
	Delete(id any) error
}
