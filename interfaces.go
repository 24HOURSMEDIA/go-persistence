package persistence

type ObjectPersisterInterface interface {
	SaveItem(key string, obj interface{}) error
	GetItem(key string, obj interface{}) error
	ListKeys() ([]string, error)
}
