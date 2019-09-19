package bolt

import (
	"errors"
	"os"

	"github.com/jeckbjy/sdb/engine"
	"github.com/jeckbjy/sdb/engine/comm"
	bolt "go.etcd.io/bbolt"
)

const name = "bolt"
const uri = "bolt.db"

func Register() {
	engine.Register(name, New)
}

func New() engine.IClient {
	return &_Client{}
}

type _Client struct {
	db    *bolt.DB
	codec engine.ICodec
}

func (c *_Client) Name() string {
	return name
}

func (c *_Client) Open(opts *engine.OpenOptions) error {
	if opts.URI == engine.DefaultURI {
		opts.URI = uri
	}

	if opts.Codec == nil {
		opts.Codec = comm.DefaultCodec()
	}

	db, err := bolt.Open(opts.URI, os.ModePerm, nil)
	if err != nil {
		return err
	}

	c.db = db
	c.codec = opts.Codec

	return nil
}

func (c *_Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}

	return nil
}

func (c *_Client) Ping() error {
	return nil
}

func (c *_Client) Database(name string) (engine.IDatabase, error) {
	if c.db == nil {
		return nil, errors.New("bad db")
	}

	db := &_Database{db: c.db, name: name, codec: c.codec}
	err := c.db.Update(func(tx *bolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists([]byte(name))
		return e
	})
	return db, err
}

func (c *_Client) Drop(name string) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(name))
	})
}
