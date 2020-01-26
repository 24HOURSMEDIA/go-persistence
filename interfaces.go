package persistence

// Object Persisters must implement this interface
type ObjectPersisterInterface interface {
	saveItemNow(key string, obj interface{}) error
	SaveItem(key string, obj interface{}) error
	GetItem(key string, obj interface{}) error
	ListKeys() ([]string, error)
}
