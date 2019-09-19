package engine

// https://docs.mongodb.com/manual/reference/operator/query-comparison/
// Query embedded document: $elemMatch vs. Dot Notation
type Token int

const (
	TOK_NULL Token = iota
	TOK_EQ
	TOK_NE
	TOK_GT
	TOK_GTE
	TOK_LT
	TOK_LTE
	TOK_IN
	TOK_NIN
	// logical Query Operators
	TOK_AND
	TOK_OR
	TOK_NOR
	TOK_NOT
)

type ICond interface {
	Operator() Token
}

// IExpr 比较操作符
type IExpr interface {
	ICond
	Key() string
	Value() interface{}
}

type IUnary interface {
	ICond
	X() ICond
}

type IBinary interface {
	ICond
	X() ICond
	Y() ICond
}

type IList interface {
	ICond
	List() []ICond
}
