package port

import "gorm.io/gorm"

type DB interface {
	Save(arg interface{}) error
	FindAll(dest interface{}, conditions map[string]interface{}) error
	List(dest interface{}, conditions map[string]interface{}, limit int, offset int) error
	FindById(dest interface{}, id string) error
	FindOne(dest interface{}, conditions map[string]interface{}) error
	FindWithPreload(dest interface{}, conditions map[string]interface{}, preload_options ...func(*gorm.DB)) error
	Delete(model interface{}, id string) error
	DeleteOne(model interface{}, conditions map[string]interface{}) error
	Update(model interface{}, condition map[string]interface{}, updates interface{}) error
	Raw(dest interface{}, query string, values ...interface{}) error
}
