package dvbcrud_go

type Repository[TModel any] interface {
    Create(model TModel) error
    Read(id any) (TModel, error)
    ReadAll() ([]TModel, error)
    Update(id any, model TModel) error
    Delete(id any) error
}
