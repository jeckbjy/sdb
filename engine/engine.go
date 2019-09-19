package engine

var gEngineMap = make(map[string]Func)

func Register(name string, fn Func) {
	gEngineMap[name] = fn
}

func New(name string) IClient {
	if name == "" {
		// return first
		for _, fn := range gEngineMap {
			return fn()
		}
	} else {
		if fn, ok := gEngineMap[name]; ok {
			return fn()
		}
	}

	return nil
}

type Func func() IClient

type IClient interface {
	Name() string
	Open(opts *OpenOptions) error
	Close() error
	Ping() error
	Database(name string) (IDatabase, error)
	Drop(name string) error
}

type IDatabase interface {
	EnsureIndex(bucket string, keys interface{}, opts *IndexOptions) (string, error)
	DropIndex(bucket string, name string) error
	Indexes(bucket string) ([]Index, error)

	// create table if not exits
	Create(bucket string, doc interface{}) error
	// delete all data and indexes from table
	Drop(bucket string) error
	// CRUD
	Insert(bucket string, doc interface{}, opts *InsertOptions) (*InsertResult, error)
	Delete(bucket string, filter ICond, opts *DeleteOptions) (*DeleteResult, error)
	Update(bucket string, filter ICond, update interface{}, opts *UpdateOptions) (*UpdateResult, error)
	Query(result interface{}, bucket string, filter ICond, opts *QueryOptions) error
}

type ICodec interface {
	Name() string
	Encode(doc interface{}) ([]byte, error)
	Decode(data []byte, doc interface{}) error
}
