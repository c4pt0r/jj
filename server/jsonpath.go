package server

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	log "github.com/ngaut/logging"
)

var re = regexp.MustCompile("^([^0-9\\s\\[][^\\s\\[]*)?(\\[-?[0-9]+\\])?$")

func isArray(v interface{}) bool {
	_, ok := v.([]interface{})
	return ok
}

func isMap(v interface{}) bool {
	_, ok := v.(map[string]interface{})
	return ok
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

func jsonPathDo(v interface{}, jp string, fn func(v interface{}), replaceWith func(v interface{}) interface{}) error {
	parts := strings.Split(jp, ".")
	for idx, token := range parts {
		sl := re.FindAllStringSubmatch(token, -1)
		if len(sl) == 0 {
			return errors.New("invalid path")
		}
		ss := sl[0]
		if ss[1] != "" {
			if dict, ok := v.(map[string]interface{}); ok {
				v = dict[ss[1]]
				if idx == len(parts)-1 && replaceWith != nil {
					newVal := replaceWith(v)
					if newVal != nil {
						dict[ss[1]] = replaceWith(v)
					}
				}
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
					if idx == len(parts)-1 && replaceWith != nil {
						newVal := replaceWith(v)
						if newVal != nil {
							array[i] = replaceWith(v)
						}
					}
				}
				if i < 0 {
					if sz+i < 0 {
						return errors.New("invalid index")
					}
					v = array[sz+i]
					if idx == len(parts)-1 && replaceWith != nil {
						newVal := replaceWith(v)
						if newVal != nil {
							array[sz+i] = replaceWith(v)
						}
					}
				}
			} else {
				return errors.New("invalid type, assume array")
			}
		}
	}
	if v != nil && fn != nil {
		fn(v)
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
			return errors.New("invalid path")
		}
	}
	return nil
}

func jsonPathQuery(v interface{}, jp string, t interface{}) error {
	return jsonPathDo(v, jp, func(v interface{}) {
		rt := reflect.ValueOf(t).Elem()
		rv := reflect.ValueOf(v)
		rt.Set(rv)
	}, nil)
}

func jsonPathIncr(v interface{}, jp string, delta int) error {
	return jsonPathDo(v, jp, nil, func(v interface{}) interface{} {
		if vv, ok := v.(float64); ok {
			return int(vv) + delta
		} else if vv, ok := v.(int); ok {
			return vv + delta
		}
		return nil
	})
}

func jsonPathPush(v interface{}, jp string, val interface{}) error {
	return jsonPathDo(v, jp, nil, func(v interface{}) interface{} {
		if vv, ok := v.([]interface{}); ok {
			return append(vv, val)
		}
		return nil
	})
}

func jsonPathPop(v interface{}, jp string, t interface{}) error {
	return jsonPathDo(v, jp, nil, func(v interface{}) interface{} {
		if vv, ok := v.([]interface{}); ok {
			rt := reflect.ValueOf(t).Elem()
			rv := reflect.ValueOf(vv[0])
			rt.Set(rv)
			return vv[1:]
		}
		return nil
	})
}

func jsonPathArrayLen(v interface{}, jp string) (int, error) {
	sz := -1
	err := jsonPathDo(v, jp, func(v interface{}) {
		if vv, ok := v.([]interface{}); ok {
			sz = len(vv)
		}
	}, nil)
	if err != nil {
		return -1, err
	}
	return sz, nil
}

func jsonPathRemove(v interface{}, jp string) error {
	return nil
}
