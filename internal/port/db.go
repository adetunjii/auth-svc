package port

type DB interface {
	Save(arg interface{}) error
	FindAll(dest interface{}, conditions map[string]interface{}) error
	List(dest interface{}, conditions map[string]interface{}, limit int, offset int) error
	FindById(dest interface{}, id string) error
	FindOne(dest interface{}, conditions map[string]interface{}) error
	FindWithPreload(dest interface{}, conditions map[string]interface{}, with string) error
	Delete(model interface{}, id string) error
	Update(model interface{}, condition map[string]interface{}, updates interface{}) error
}
