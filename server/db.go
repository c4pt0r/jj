package server

import (
	"errors"
	"fmt"
)

var ErrNoSuchKey = errors.New("no such key")

type Db interface {
	PutDoc(key string, val interface{}) error
	GetDoc(key string) (interface{}, error)
	RemoveDoc(key string) error
	PutPath(key string, path string, val interface{}) error
	GetPath(key string, path string) (interface{}, error)
	IncrPath(key string, path string, val interface{}) error
	PushPath(key string, path string, val interface{}) error
	PopPath(key string, path string) (interface{}, error)
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
	m map[string]interface{}
}

func NewMapDb() *MapDb {
	return &MapDb{
		m: make(map[string]interface{}),
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
	if val, ok := db.m[key]; ok {
		var ret interface{}
		err := jsonPathQuery(val, path, &ret)
		if err != nil {
			return nil, err
		}
		return ret, nil
	}
	return nil, ErrNoSuchKey
}

func (db *MapDb) PutPath(key string, path string, val interface{}) error {
	if v, ok := db.m[key]; ok {
		return jsonPathSet(v, path, val)
	}
	return ErrNoSuchKey
}

func (db *MapDb) IncrPath(key string, path string, val interface{}) error {
	if v, ok := db.m[key]; ok {
		if delta, ok := val.(float64); ok {
			return jsonPathIncr(v, path, int(delta))
		} else {
			fmt.Errorf("type error %v", val)
		}
	}
	return ErrNoSuchKey
}

func (db *MapDb) PushPath(key string, path string, val interface{}) error {
	if v, ok := db.m[key]; ok {
		return jsonPathPush(v, path, val)
	}
	return ErrNoSuchKey
}

func (db *MapDb) PopPath(key string, path string) (interface{}, error) {
	if val, ok := db.m[key]; ok {
		var ret interface{}
		err := jsonPathPop(val, path, &ret)
		if err != nil {
			return nil, err
		}
		return ret, nil
	}
	return nil, ErrNoSuchKey
}

func (db *MapDb) RemovePath(key string, path string) error {
	// TODO
	return fmt.Errorf("not implement yet")
}

func (db *MapDb) Scan(keyPrefix string) (KVIter, error) {
	// TODO
	return nil, fmt.Errorf("not implement yet")
}

func (db *MapDb) Save(fileName string, context interface{}) error {
	// TODO
	return fmt.Errorf("not implement yet")
}
