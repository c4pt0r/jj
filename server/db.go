package server

type Db interface {
	PutDoc(key string, val interface{}) error
	GetDoc(key string) (interface{}, error)
	RemoveDoc(key string) error
	PutPath(key string, path string, val interface{}) error
	GetPath(key string, path string) (interface{}, error)
	RemovePath(key string, path string) error
	Scan(keyPrefix string) (KVIter, error)
	Save(fileName string, context interface{}) error
}

type KVIter interface {
	Next() (KVIter, error)
	HasNext() bool
	Val() string
	Key() interface{}
}

type MapDb struct {
	m     map[string]interface{}
	dirty int
}

func NewMapDb() *MapDb {
	return &MapDb{
		m:     make(map[string]interface{}),
		dirty: 0,
	}
}

func (db *MapDb) PutDoc(key string, val interface{}) error {
	db.m[key] = val
	return nil
}

func (db *MapDb) GetDoc(key string) (interface{}, error) {
	val := db.m[key]
	return val, nil
}

func (db *MapDb) RemoveDoc(key string) error {
	delete(db.m, key)
	return nil
}

func (db *MapDb) GetPath(key string, path string) (interface{}, error) {
	return nil, nil
}

func (db *MapDb) PutPath(key string, path string, val interface{}) error {
	return nil
}

func (db *MapDb) RemovePath(key string, path string) error {
	return nil
}

func (db *MapDb) Scan(keyPrefix string) KVIter {
	return nil
}
