package sql

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/jeckbjy/sdb/engine"
)

func toWhere(cond engine.ICond) (string, error) {
	// 全部
	if cond == nil {
		return "", nil
	}

	// 条件
	b := whereBuilder{}
	b.Write("WHERE")
	if err := b.build(cond, 0); err != nil {
		return "", err
	}

	return b.String(), nil
}

func toValue(v interface{}) string {
	return toString(reflect.ValueOf(v))
}

func toString(v reflect.Value) string {
	switch kind := v.Type().Kind(); {
	case kind == reflect.String:
		return fmt.Sprintf("'%s'", v)
	case kind >= reflect.Bool && kind <= reflect.Float64 && kind != reflect.Uintptr:
		return fmt.Sprintf("%+v", v.Interface())
	case kind == reflect.Slice:
		str := make([]string, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			str = append(str, toString(v.Index(i)))
		}
		return strings.Join(str, ",")
	case kind == reflect.Ptr:
		return toString(v.Elem())
	default:
		return ""
	}
}

var operators = [...]string{
	engine.TOK_EQ:  "=",
	engine.TOK_NE:  "<>",
	engine.TOK_GT:  ">",
	engine.TOK_GTE: ">=",
	engine.TOK_LT:  "<",
	engine.TOK_LTE: "<=",
	engine.TOK_IN:  "IN",
	engine.TOK_NIN: "NOT IN",
	engine.TOK_AND: "AND",
	engine.TOK_OR:  "OR",
	engine.TOK_NOT: "NOT",
}

type whereBuilder struct {
	builder
}

func (b *whereBuilder) build(cond engine.ICond, depth int) error {
	switch op := cond.Operator(); {
	case op >= engine.TOK_EQ && op <= engine.TOK_LTE:
		s := cond.(engine.IExpr)
		// assert s.Value?
		b.Writef("%s %s %s", s.Key(), operators[op], toValue(s.Value()))
	case op == engine.TOK_IN || op == engine.TOK_NIN:
		s := cond.(engine.IExpr)
		b.Writef("%s %s (%s)", s.Key(), operators[op], toValue(s.Value()))
	case op == engine.TOK_AND || op == engine.TOK_OR:
		s := cond.(engine.IList)
		if len(s.List()) == 0 {
			return nil
		}
		if depth > 0 {
			b.Write("(")
		}
		for i, l := range s.List() {
			if i > 0 {
				b.Write(operators[op])
			}
			if err := b.build(l, depth+1); err != nil {
				return err
			}
		}

		if depth > 0 {
			b.Write(")")
		}
	default:
		return errors.New("not support")
	}

	return nil
}
