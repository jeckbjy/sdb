package bolt

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/jeckbjy/sdb/engine"
	"github.com/jeckbjy/sdb/engine/comm"
	bolt "go.etcd.io/bbolt"
)

var errBadBucket = errors.New("bad bucket")

// 目前不支持索引,更新删除无脑遍历
type _Database struct {
	db    *bolt.DB
	codec engine.ICodec
	name  string
}

func (d *_Database) getBucket(tx *bolt.Tx, name string) (*bolt.Bucket, error) {
	db := tx.Bucket([]byte(d.name))
	if db == nil {
		return nil, errBadBucket
	}

	b, err := db.CreateBucketIfNotExists([]byte(name))
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (d *_Database) EnsureIndex(bucket string, keys interface{}, opts *engine.IndexOptions) (string, error) {
	return "", nil
}

func (d *_Database) DropIndex(bucket string, name string) error {
	return nil
}

func (d *_Database) Indexes(bucket string) ([]engine.Index, error) {
	return nil, nil
}

func (d *_Database) Create(bucket string, doc interface{}) error {
	return nil
}

func (d *_Database) Drop(bucket string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		db := tx.Bucket([]byte(d.name))
		if db == nil {
			return errBadBucket
		}

		return db.DeleteBucket([]byte(bucket))
	})
}

func (d *_Database) Insert(bucket string, doc interface{}, opts *engine.InsertOptions) (*engine.InsertResult, error) {
	res := &engine.InsertResult{}
	err := d.db.Update(func(tx *bolt.Tx) error {
		b, err := d.getBucket(tx, bucket)
		if err != nil {
			return err
		}

		if opts.One {
			id, err := d.put(b, doc, true)
			if err != nil {
				return err
			}
			res.InsertedIDs = append(res.InsertedIDs, id)
		} else {
			v := reflect.ValueOf(doc)
			if v.Type().Kind() != reflect.Slice {
				return errors.New("insert fail,must slice")
			}

			for i := 0; i < v.Len(); i++ {
				id, err := d.put(b, v.Index(i).Interface(), true)
				if err != nil {
					return err
				}
				res.InsertedIDs = append(res.InsertedIDs, id)
			}
		}

		return nil
	})

	return res, err
}

func (d *_Database) put(b *bolt.Bucket, doc interface{}, insert bool) (interface{}, error) {
	// index?
	id, err := comm.GetID(doc)
	if err != nil {
		return nil, err
	}

	// no id, auto create?
	sid := fmt.Sprintf("%+v", id)
	if sid == "" {
		sid = comm.NewID()
		id = sid
		m := comm.ToMap(d.codec, doc)
		m["_id"] = sid
		doc = m
	}

	if insert && b.Get([]byte(sid)) != nil {
		return id, fmt.Errorf("repeated insertId,%+v", sid)
	}

	data, err := d.codec.Encode(doc)
	if err != nil {
		return nil, err
	}

	return id, b.Put([]byte(sid), data)
}

func (d *_Database) Delete(bucket string, filter engine.ICond, opts *engine.DeleteOptions) (*engine.DeleteResult, error) {
	res := &engine.DeleteResult{}
	err := d.db.Update(func(tx *bolt.Tx) error {
		// 遍历
		b, err := d.getBucket(tx, bucket)
		if err != nil {
			return err
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			// decode
			doc := make(map[string]interface{})
			if err := d.codec.Decode(v, &doc); err != nil {
				continue
			}

			if comm.Apply(filter, doc) {
				if b.Delete(k) == nil {
					res.DeletedCount++
				}
			}
		}
		return nil
	})

	return res, err
}

func (d *_Database) Update(bucket string, filter engine.ICond, update interface{}, opts *engine.UpdateOptions) (*engine.UpdateResult, error) {
	res := &engine.UpdateResult{}
	err := d.db.Update(func(tx *bolt.Tx) error {
		b, err := d.getBucket(tx, bucket)
		if err != nil {
			return err
		}

		mm := comm.ToMap(d.codec, update)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			// decode
			doc := make(map[string]interface{})
			if err := d.codec.Decode(v, &doc); err != nil {
				continue
			}

			if comm.Apply(filter, doc) {
				res.MatchedCount++
				// merge and save
				if err := comm.Merge(doc, mm); err == nil {
					if _, err = d.put(b, doc, false); err != nil {
						return err
					}
				} else {
					return err
				}

				if opts.One {
					break
				}
			}
		}

		if res.MatchedCount == 0 && opts.Upsert {
			_, err := d.put(b, mm, true)
			return err
		}

		return nil
	})

	return res, err
}

func (d *_Database) Query(result interface{}, bucket string, filter engine.ICond, opts *engine.QueryOptions) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b, err := d.getBucket(tx, bucket)
		if err != nil {
			return err
		}

		c := b.Cursor()
		if opts.One {
			for k, v := c.First(); k != nil; k, v = c.Next() {
				doc := make(map[string]interface{})
				if err := d.codec.Decode(v, &doc); err != nil {
					continue
				}

				if comm.Apply(filter, doc) {
					return d.codec.Decode(v, result)
				}
			}
		} else {
			finder := comm.Finder{}
			for k, v := c.First(); k != nil; k, v = c.Next() {
				doc := make(map[string]interface{})
				if err := d.codec.Decode(v, &doc); err != nil {
					continue
				}

				if comm.Apply(filter, doc) {
					finder.Push(doc, v)
				}
			}

			finder.Sort(opts.Sort)
			return finder.Find(result, d.codec, opts.Skip, opts.Limit, opts.Projection)
		}

		return nil
	})
}
