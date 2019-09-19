package comm

import (
	"reflect"
	"time"
)

func cmpEqual(x interface{}, y interface{}) bool {
	if x == nil || y == nil {
		return x == nil && y == nil
	}

	//
	k1 := reflect.TypeOf(x).Kind()
	k2 := reflect.TypeOf(y).Kind()
	if k1 != k2 {
		return false //?
	}

	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)

	switch k1 {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v1.Int() == v2.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v1.Uint() == v2.Uint()
	case reflect.Float32, reflect.Float64:
		return v1.Float() == v2.Float()
	case reflect.String:
		return v1.String() == v2.String()
	case reflect.Struct:
		t1, ok1 := v1.Interface().(time.Time)
		t2, ok2 := v2.Interface().(time.Time)
		if ok1 && ok2 && t1 == t2 {
			return true
		}
	default:
		return false
	}
	return false
}

// x > y return true,otherwise return false
func cmpGreat(x interface{}, y interface{}) bool {
	if y == nil {
		return true
	}
	//
	k1 := reflect.TypeOf(x).Kind()
	k2 := reflect.TypeOf(y).Kind()
	if k1 != k2 {
		return false //?
	}

	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)

	switch k1 {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v1.Int() > v2.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v1.Uint() > v2.Uint()
	case reflect.Float32, reflect.Float64:
		return v1.Float() > v2.Float()
	case reflect.String:
		return v1.String() > v2.String()
	case reflect.Struct:
		t1, ok1 := v1.Interface().(time.Time)
		t2, ok2 := v2.Interface().(time.Time)
		if ok1 && ok2 && t1.After(t2) {
			return true
		}
	default:
		return false
	}
	return false
}

// x < y return true,otherwise return false
func cmpLess(x interface{}, y interface{}) bool {
	if y == nil {
		return true
	}
	//
	k1 := reflect.TypeOf(x).Kind()
	k2 := reflect.TypeOf(y).Kind()
	if k1 != k2 {
		return false //?
	}

	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)

	switch k1 {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v1.Int() < v2.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v1.Uint() < v2.Uint()
	case reflect.Float32, reflect.Float64:
		return v1.Float() < v2.Float()
	case reflect.String:
		return v1.String() < v2.String()
	case reflect.Struct:
		t1, ok1 := v1.Interface().(time.Time)
		t2, ok2 := v2.Interface().(time.Time)
		if ok1 && ok2 && t1.Before(t2) {
			return true
		}
	default:
		return false
	}
	return false
}

func cmpIn(x interface{}, y interface{}) bool {
	if reflect.TypeOf(y).Kind() != reflect.Slice {
		return cmpEqual(x, y)
	}

	s := reflect.ValueOf(y)
	for i := 0; i < s.Len(); i++ {
		f := s.Index(i)
		if cmpEqual(x, f.Interface()) {
			return true
		}
	}

	return false
}
