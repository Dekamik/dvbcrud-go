package dvbcrud_go

type Repository[TId any, TModel any] interface {
    Create(model TModel) (TModel, error)
    Read(id TId) (TModel, error)
    ReadAll() ([]TModel, error)
    Update(id TId, model TModel) (TModel, error)
    Delete(id TId) error
}

type repository[TId any, TModel any] struct {
    tableName string
}

func (r repository[TId, TModel]) Create(model TModel) (TModel, error) {
    return nil, nil
}

func (r repository[TId, TModel]) Read(id TId) (TModel, error) {
    return nil, nil
}

func (r repository[TId, TModel]) ReadAll() ([]TModel, error) {
    return nil, nil
}

func (r repository[TId, TModel]) Update(id TId, model TModel) (TModel, error) {
    return nil, nil
}

func (r repository[TId, TModel]) Delete(id TId) error {
    return nil
}

func New[TId any, TModel any](tableName string) Repository[TId, TModel] {
    return repository[TId, TModel]{
        tableName: tableName,
    }
}
