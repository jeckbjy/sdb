package sdb

import "github.com/jeckbjy/sdb/engine"

type OpenOption func(*engine.OpenOptions)
type IndexOption func(*engine.IndexOptions)
type InsertOption func(*engine.InsertOptions)
type DeleteOption func(*engine.DeleteOptions)
type UpdateOption func(*engine.UpdateOptions)
type QueryOption func(*engine.QueryOptions)

const DefaultURI = engine.DefaultURI

func WithDriver(d string) OpenOption {
	return func(o *engine.OpenOptions) {
		o.Driver = d
	}
}

func WithURI(uri string) OpenOption {
	return func(o *engine.OpenOptions) {
		o.URI = uri
	}
}

func WithClient(client interface{}) OpenOption {
	return func(o *engine.OpenOptions) {
		o.Client = client
	}
}

func WithCodec(codec engine.ICodec) OpenOption {
	return func(o *engine.OpenOptions) {
		o.Codec = codec
	}
}

func WithUpsert() UpdateOption {
	return func(o *engine.UpdateOptions) {
		o.Upsert = true
	}
}

func WithSkip(skip int64) QueryOption {
	return func(o *engine.QueryOptions) {
		o.Skip = skip
	}
}

func WithLimit(limit int64) QueryOption {
	return func(o *engine.QueryOptions) {
		o.Limit = limit
	}
}

func WithSort(key string, asc bool) QueryOption {
	return func(o *engine.QueryOptions) {
		if o.Sort == nil {
			o.Sort = make(map[string]int)
		}
		if asc {
			o.Sort[key] = 1
		} else {
			o.Sort[key] = -1
		}
	}
}

func WithProjection(fields []string, including bool) QueryOption {
	return func(o *engine.QueryOptions) {
		if o.Projection == nil {
			o.Projection = make(map[string]int)
		}
		v := 0
		if including {
			v = 1
		}
		for _, k := range fields {
			o.Projection[k] = v
		}
	}
}
