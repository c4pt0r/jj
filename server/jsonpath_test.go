package server

import (
	"encoding/json"
	"log"
	"testing"
)

var (
	doc1 = `
	{
		"a" : 1.1,
		"b" : [1,2,3,4,5],
		"c" : [{"d":[1,2,3,4,5]}, {"e": 2}]
	}
`

	map1 = map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": []int{1, 2, 3, 4, 5},
	}
)

func TestJsonPathSet(t *testing.T) {
	var v interface{}
	err := json.Unmarshal([]byte(doc1), &v)
	if err != nil {
		t.Error(err)
	}

	var val interface{}
	jsonPathQuery(v, "a", &val)
	log.Println(val.(float64))

	err = jsonPathIncr(v, "a", 1)
	if err != nil {
		t.Error(err)
	}

	jsonPathQuery(v, "a", &val)
	log.Println(val)

	err = jsonPathIncr(v, "b[20]", 100)
	if err == nil {
		t.Error("should error")
	}

	jsonPathQuery(v, "b[20]", &val)
	log.Println(val)

	jsonPathPush(v, "b", 6)
	jsonPathQuery(v, "b", &val)
	log.Println(val)

	var i interface{}
	jsonPathPop(v, "b", &i)
	log.Println(i)
	jsonPathQuery(v, "b", &val)
	log.Println(val)

	jsonPathRemove(v, "c")
	jsonPathQuery(v, ".", &val)
	log.Println(val)
}
