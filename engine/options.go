package engine

const DefaultURI = "*"

type OpenOptions struct {
	Client interface{} // 原生的Client
	Driver string
	URI    string
	Codec  ICodec
}

type IndexOptions struct {
	Name       string
	Background bool
	Sparse     bool
	Unique     bool
}

type InsertOptions struct {
	One bool
}

type DeleteOptions struct {
	One bool
}

type UpdateOptions struct {
	One    bool
	Upsert bool
}

type QueryOptions struct {
	One        bool
	Skip       int64
	Limit      int64
	Sort       map[string]int //1:ascending, -1:descending
	Projection map[string]int //1:include 0:exclude(有些不支持),nil代表全部
}
