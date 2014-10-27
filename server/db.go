package server

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	log "github.com/ngaut/logging"
)

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

var re = regexp.MustCompile("^([^0-9\\s\\[][^\\s\\[]*)?(\\[-?[0-9]+\\])?$")

func isArray(v interface{}) bool {
	if _, ok := v.([]interface{}); ok {
		return true
	}
	return false
}

func isMap(v interface{}) bool {
	if _, ok := v.(map[string]interface{}); ok {
		return true
	}
	return false
}

func jsonPathQuery(v interface{}, jp string, t interface{}) error {
	for _, token := range strings.Split(jp, ".") {
		sl := re.FindAllStringSubmatch(token, -1)
		if len(sl) == 0 {
			return errors.New("invalid path")
		}
		ss := sl[0]
		if ss[1] != "" {
			if dict, ok := v.(map[string]interface{}); ok {
				v = dict[ss[1]]
			} else {
				return errors.New("invalid type, assume dict")
			}
		}
		if ss[2] != "" {
			i, err := strconv.Atoi(ss[2][1 : len(ss[2])-1])
			if err != nil {
				return err
			}
			if array, ok := v.([]interface{}); ok {
				sz := len(array)
				if i >= 0 {
					if i > sz-1 {
						return errors.New("invalid index")
					}
					v = array[i]
				}
				if i < 0 {
					if sz+i < 0 {
						return errors.New("invalid index")
					}
					v = array[sz+i]
				}
			} else {
				return errors.New("invalid type, assume array")
			}
		}
	}
	rt := reflect.ValueOf(t).Elem()
	rv := reflect.ValueOf(v)
	rt.Set(rv)
	return nil
}

func getItemFromMap(o interface{}, k string) interface{} {
	if m, ok := o.(map[string]interface{}); ok {
		return m[k]
	} else {
		return nil
	}
}

func getItemFromArray(o interface{}, i int) interface{} {
	if a, ok := o.([]interface{}); ok {
		if i >= 0 && i < len(a) {
			return a[i]
		}
	}
	return nil
}

func jsonPathSet(v interface{}, jp string, val interface{}) error {
	cur := v
	parts := strings.Split(jp, ".")
	for idx, part := range parts {
		sl := re.FindAllStringSubmatch(part, -1)
		if len(sl) == 0 {
			return errors.New("invalid path")
		}
		ss := sl[0]

		log.Info(ss[1], ss[2], cur)
		if len(ss[1]) > 0 && len(ss[2]) == 0 && isMap(cur) {
			// handle x
			if idx != len(parts)-1 {
				i := getItemFromMap(cur, ss[1])
				if i == nil || (!isArray(i) && !isMap(i)) {
					cur.(map[string]interface{})[ss[1]] = make(map[string]interface{})
				}
				cur = cur.(map[string]interface{})[ss[1]]
			} else {
				cur.(map[string]interface{})[ss[1]] = val
			}
		} else if len(ss[1]) == 0 && len(ss[2]) > 0 && isArray(cur) {
			// handle [i]
			i, err := strconv.Atoi(ss[2][1 : len(ss[2])-1])
			if err != nil {
				return errors.New("invalid path")
			}
			if idx != len(parts)-1 {
				v := getItemFromArray(cur, i)
				if v == nil {
					return errors.New("invalid path")
				}
				cur = v
			} else {
				cur.([]interface{})[i] = val
			}
		} else if len(ss[1]) > 0 && len(ss[2]) > 0 && isMap(cur) && isArray(getItemFromMap(cur, ss[1])) {
			// handle x[i]
			i, err := strconv.Atoi(ss[2][1 : len(ss[2])-1])
			if err != nil {
				return errors.New("invalid path")
			}
			if idx != len(parts)-1 {
				v := getItemFromArray(getItemFromMap(cur, ss[1]), i)
				if v == nil {
					return errors.New("invalid path")
				}
				cur = v
			} else {
				v := getItemFromMap(cur, ss[1])
				if i >= 0 && i < len(v.([]interface{})) {
					v.([]interface{})[i] = val
				}
			}
		} else {
			log.Warning("invalid path")
			return errors.New("invalid path")
		}
	}
	return nil
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
	return nil, fmt.Errorf("no such key: %s", key)
}

func (db *MapDb) PutPath(key string, path string, val interface{}) error {
	var v interface{}
	if _, ok := db.m[key]; !ok {
		return fmt.Errorf("no such key: %s", key)
	}
	v = db.m[key]
	return jsonPathSet(v, path, val)
}

func (db *MapDb) RemovePath(key string, path string) error {
	return nil
}

func (db *MapDb) Scan(keyPrefix string) (KVIter, error) {
	return nil, nil
}

func (db *MapDb) Save(fileName string, context interface{}) error {
	return nil
}
