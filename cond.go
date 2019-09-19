package sdb

import "github.com/jeckbjy/sdb/engine"

type expr struct {
	op  engine.Token
	key string
	val interface{}
}

func (e *expr) Operator() engine.Token {
	return e.op
}

func (e *expr) Key() string {
	return e.key
}

func (e *expr) Value() interface{} {
	return e.val
}

type unary struct {
	op engine.Token
	x  engine.ICond
}

func (u *unary) Operator() engine.Token {
	return u.op
}

func (u *unary) X() engine.ICond {
	return u.x
}

type binary struct {
	op engine.Token
	x  engine.ICond
	y  engine.ICond
}

func (b *binary) Operator() engine.Token {
	return b.op
}

func (b *binary) X() engine.ICond {
	return b.x
}

func (b *binary) Y() engine.ICond {
	return b.y
}

type list struct {
	op    engine.Token
	conds []engine.ICond
}

func (l *list) Operator() engine.Token {
	return l.op
}

func (l *list) List() []engine.ICond {
	return l.conds
}

func Eq(key string, value interface{}) engine.ICond {
	return &expr{op: engine.TOK_EQ, key: key, val: value}
}

func Gt(key string, value interface{}) engine.ICond {
	return &expr{engine.TOK_GT, key, value}
}

func Gte(key string, value interface{}) engine.ICond {
	return &expr{engine.TOK_GTE, key, value}
}

func Lt(key string, value interface{}) engine.ICond {
	return &expr{engine.TOK_LT, key, value}
}

func Lte(key string, value interface{}) engine.ICond {
	return &expr{engine.TOK_LTE, key, value}
}

func Ne(key string, value interface{}) engine.ICond {
	return &expr{engine.TOK_NE, key, value}
}

func In(key string, value interface{}) engine.ICond {
	return &expr{engine.TOK_IN, key, value}
}

func Nin(key string, value interface{}) engine.ICond {
	return &expr{engine.TOK_NIN, key, value}
}

func Not(cond engine.ICond) engine.ICond {
	return &unary{engine.TOK_NOT, cond}
}

func And(conds ...engine.ICond) engine.ICond {
	return &list{engine.TOK_AND, conds}
}

func Or(conds ...engine.ICond) engine.ICond {
	return &list{engine.TOK_OR, conds}
}

func Nor(conds ...engine.ICond) engine.ICond {
	return &list{engine.TOK_NOR, conds}
}
