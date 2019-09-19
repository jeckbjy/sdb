package sdb

import (
	"errors"
	"reflect"

	"github.com/jeckbjy/sdb/engine"
)

type Database struct {
	db engine.IDatabase
}

func (d *Database) EnsureIndex(bucket string, keys interface{}, opts ...IndexOption) (string, error) {
	o := engine.IndexOptions{}
	for _, fn := range opts {
		fn(&o)
	}

	return d.db.EnsureIndex(bucket, keys, &o)
}

func (d *Database) DropIndex(bucket string, name string) error {
	return d.db.DropIndex(bucket, name)
}

func (d *Database) Indexes(bucket string) ([]engine.Index, error) {
	return d.db.Indexes(bucket)
}

func (d *Database) Insert(bucket string, doc interface{}, opts ...InsertOption) (*engine.InsertResult, error) {
	o := engine.InsertOptions{}
	o.One = reflect.TypeOf(doc).Kind() != reflect.Slice
	for _, fn := range opts {
		fn(&o)
	}

	return d.db.Insert(bucket, doc, &o)
}

func (d *Database) Delete(bucket string, filter engine.ICond, opts ...DeleteOption) (*engine.DeleteResult, error) {
	o := engine.DeleteOptions{}
	o.One = false
	for _, fn := range opts {
		fn(&o)
	}

	return d.db.Delete(bucket, filter, &o)
}

func (d *Database) Update(bucket string, filter engine.ICond, update interface{}, opts ...UpdateOption) (*engine.UpdateResult, error) {
	o := engine.UpdateOptions{}
	o.One = reflect.TypeOf(update).Kind() != reflect.Slice
	for _, fn := range opts {
		fn(&o)
	}

	return d.db.Update(bucket, filter, update, &o)
}

func (d *Database) Query(result interface{}, bucket string, filter engine.ICond, opts ...QueryOption) error {
	typ := reflect.TypeOf(result)
	if typ.Kind() != reflect.Ptr {
		return errors.New("find result must be ptr")
	}

	o := engine.QueryOptions{}
	o.One = typ.Elem().Kind() != reflect.Slice
	for _, fn := range opts {
		fn(&o)
	}

	return d.db.Query(result, bucket, filter, &o)
}

func (d *Database) Create(bucket string, doc interface{}) error {
	return d.db.Create(bucket, doc)
}

func (d *Database) Drop(bucket string) error {
	return d.db.Drop(bucket)
}
