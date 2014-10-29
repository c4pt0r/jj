package server

import (
	"errors"
	"fmt"
	"hash/crc32"
	"sync"
)

var ErrNoSuchKey = errors.New("no such key")

const (
	MaxSlotSize = 1024
)

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

type Slot struct {
	m    map[string]interface{}
	lock sync.RWMutex
}

func NewSlot() *Slot {
	return &Slot{
		m:    make(map[string]interface{}),
		lock: sync.RWMutex{},
	}
}

type MapDb struct {
	slots    []*Slot
	keyCount int
}

func NewMapDb() *MapDb {
	// init slots
	var slots []*Slot
	for i := 0; i < MaxSlotSize; i++ {
		slots = append(slots, NewSlot())
	}
	return &MapDb{
		slots:    slots,
		keyCount: 0,
	}
}

func GetSlotIdFromKey(key string) int {
	h := crc32.Checksum([]byte(key), crc32.IEEETable)
	return int(h) % MaxSlotSize
}

func (db *MapDb) PutDoc(key string, val interface{}) error {
	id := GetSlotIdFromKey(key)
	db.slots[id].lock.Lock()
	db.slots[id].m[key] = val
	db.slots[id].lock.Unlock()
	return nil
}

func (db *MapDb) GetDoc(key string) (interface{}, error) {
	id := GetSlotIdFromKey(key)
	db.slots[id].lock.RLock()
	val := db.slots[id].m[key]
	db.slots[id].lock.RUnlock()
	return val, nil
}

func (db *MapDb) RemoveDoc(key string) error {
	id := GetSlotIdFromKey(key)
	db.slots[id].lock.Lock()
	delete(db.slots[id].m, key)
	db.slots[id].lock.Unlock()
	return nil
}

func (db *MapDb) GetPath(key string, path string) (interface{}, error) {
	id := GetSlotIdFromKey(key)
	db.slots[id].lock.RLock()
	val, ok := db.slots[id].m[key]
	db.slots[id].lock.RUnlock()

	if ok {
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
	id := GetSlotIdFromKey(key)
	db.slots[id].lock.Lock()
	defer db.slots[id].lock.Unlock()
	if v, ok := db.slots[id].m[key]; ok {
		return jsonPathSet(v, path, val)
	}
	return ErrNoSuchKey
}

func (db *MapDb) IncrPath(key string, path string, val interface{}) error {
	id := GetSlotIdFromKey(key)
	db.slots[id].lock.Lock()
	defer db.slots[id].lock.Unlock()
	if v, ok := db.slots[id].m[key]; ok {
		if delta, ok := val.(float64); ok {
			return jsonPathIncr(v, path, int(delta))
		} else {
			fmt.Errorf("type error %v", val)
		}
	}
	return ErrNoSuchKey
}

func (db *MapDb) PushPath(key string, path string, val interface{}) error {
	id := GetSlotIdFromKey(key)
	db.slots[id].lock.Lock()
	defer db.slots[id].lock.Unlock()
	if v, ok := db.slots[id].m[key]; ok {
		return jsonPathPush(v, path, val)
	}
	return ErrNoSuchKey
}

func (db *MapDb) PopPath(key string, path string) (interface{}, error) {
	id := GetSlotIdFromKey(key)
	db.slots[id].lock.Lock()
	defer db.slots[id].lock.Unlock()
	if v, ok := db.slots[id].m[key]; ok {
		var ret interface{}
		err := jsonPathPop(v, path, &ret)
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
