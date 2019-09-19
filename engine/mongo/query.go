package mongo

import (
	"errors"

	"github.com/jeckbjy/sdb/engine"
	"go.mongodb.org/mongo-driver/bson"
)

// https://docs.mongodb.com/manual/reference/operator/query-comparison/
func toBson(cond engine.ICond) (bson.M, error) {
	if cond == nil {
		return bson.M{}, nil
	}

	m := make(bson.M)
	e := build(m, cond)
	return m, e
}

func build(b bson.M, cond engine.ICond) error {
	op := operators[cond.Operator()]
	switch cond.(type) {
	case engine.IExpr:
		s := cond.(engine.IExpr)
		b[s.Key()] = bson.M{op: s.Value()}
		return nil
	case engine.IUnary:
		s := cond.(engine.IUnary)
		n, err := toBson(s.X())
		if err != nil {
			return err
		}
		b[op] = n
	case engine.IList:
		s := cond.(engine.IList)
		conds := make([]interface{}, 0, len(s.List()))
		for _, l := range s.List() {
			n, e := toBson(l)
			if e != nil {
				return e
			}
			conds = append(conds, n)
		}
		b[op] = conds
		return nil
	default:
		return errors.New("not support")
	}
	return nil
}

var operators = [...]string{
	engine.TOK_EQ:  "$eq",
	engine.TOK_NE:  "$ne",
	engine.TOK_GT:  "$gt",
	engine.TOK_GTE: "$gte",
	engine.TOK_LT:  "$lt",
	engine.TOK_LTE: "$lte",
	engine.TOK_IN:  "$in",
	engine.TOK_NIN: "$nin",
	engine.TOK_AND: "$and",
	engine.TOK_OR:  "$or",
	engine.TOK_NOR: "$nor",
	engine.TOK_NOT: "$not",
}
