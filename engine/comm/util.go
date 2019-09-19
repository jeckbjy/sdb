package comm

import (
	"errors"
	"reflect"
	"strings"

	"github.com/jeckbjy/sdb/engine"
	"github.com/jeckbjy/sdb/engine/comm/xid"
)

// 自动创建ID
func NewID() string {
	return xid.New().String()
}

// 查找唯一ID,_id
func GetID(doc interface{}) (interface{}, error) {
	if m, ok := doc.(map[string]interface{}); ok {
		return m["_id"], nil
	}

	t := reflect.TypeOf(doc)
	v := reflect.ValueOf(doc)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, errors.New("bad doc")
	}

	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("sdb")
		if tag == "_id" {
			return v.Field(i).Interface(), nil
		}

		if strings.ToLower(f.Name) == "id" {
			return v.Field(i).Interface(), nil
		}
	}

	return nil, nil
}

// 校验是否符合条件
func Apply(cond engine.ICond, doc map[string]interface{}) bool {
	if cond == nil {
		return true
	}

	switch cond.(type) {
	case engine.IExpr:
		s := cond.(engine.IExpr)
		x := GetField(doc, s.Key())
		y := s.Value()
		switch cond.Operator() {
		case engine.TOK_EQ:
			return cmpEqual(x, y)
		case engine.TOK_NE:
			return !cmpEqual(x, y)
		case engine.TOK_GT:
			return cmpGreat(x, y)
		case engine.TOK_GTE:
			return !cmpLess(y, x)
		case engine.TOK_LT:
			return cmpLess(x, y)
		case engine.TOK_LTE:
			return !cmpGreat(y, x)
		case engine.TOK_IN:
			return cmpIn(x, y)
		case engine.TOK_NIN:
			return !cmpIn(x, y)
		}
	case engine.IUnary:
		s := cond.(engine.IUnary)
		switch cond.Operator() {
		case engine.TOK_NOT:
			return !Apply(s.X(), doc)
		}
	case engine.IList:
		s := cond.(engine.IList)
		switch cond.Operator() {
		case engine.TOK_AND:
			for _, c := range s.List() {
				if !Apply(c, doc) {
					return false
				}
			}
			return true
		case engine.TOK_OR:
			// a | b | c
			for _, c := range s.List() {
				if Apply(c, doc) {
					return true
				}
			}

			return false
		case engine.TOK_NOR:
			// !(a|b|c) == !a && !b && !c
			for _, c := range s.List() {
				if Apply(c, doc) {
					return false
				}
			}
			return true
		}
	}

	return false
}

func GetField(doc map[string]interface{}, key string) interface{} {
	tokens := strings.Split(key, ".")
	if len(tokens) <= 1 {
		return doc[key]
	}

	m := doc
	for i, k := range tokens {
		v, ok := m[k]
		if !ok {
			return nil
		}

		if i == len(tokens)-1 {
			return v
		}

		if n, ok := v.(map[string]interface{}); !ok {
			return nil
		} else {
			m = n
		}
	}

	return nil
}

func ToMap(codec engine.ICodec, doc interface{}) map[string]interface{} {
	if m, ok := doc.(map[string]interface{}); ok {
		return m
	} else {
		d, err := codec.Encode(doc)
		if err == nil {
			m := make(map[string]interface{})
			if codec.Decode(d, m) == nil {
				return m
			}
		}

		return nil
	}
}

// 把b中数据合并到a中
func Merge(a map[string]interface{}, b map[string]interface{}) error {
	for k, v := range b {
		tokens := strings.Split(k, ".")
		if len(tokens) > 1 {
			if err := setField(a, tokens, v); err != nil {
				return err
			}
		} else {
			a[k] = v
		}
	}

	return nil
}

func setField(m map[string]interface{}, tokens []string, field interface{}) error {
	key := tokens[0]
	if len(tokens) == 1 {
		m[key] = field
		return nil
	}

	if s, ok := m[key]; !ok {
		n := make(map[string]interface{})
		m[key] = n
		return setField(n, tokens[1:], field)
	} else if n, ok := s.(map[string]interface{}); ok {
		return setField(n, tokens[1:], field)
	} else {
		return errors.New("bad format")
	}
}
