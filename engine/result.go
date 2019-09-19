package engine

type InsertResult struct {
	InsertedIDs []interface{}
}

type DeleteResult struct {
	DeletedCount int64
}

// UpdateResult is a result of an update operation.
//
// UpsertedID will be a Go type that corresponds to a BSON type.
type UpdateResult struct {
	// The number of documents that matched the filter.
	MatchedCount int64
	// The number of documents that were modified.
	ModifiedCount int64
	// The number of documents that were upserted.
	UpsertedCount int64
	// The identifier of the inserted document if an upsert took place.
	UpsertedID interface{}
}

type Index struct {
	Name       string
	Keys       []string
	Background bool
	Sparse     bool
	Unique     bool
}
