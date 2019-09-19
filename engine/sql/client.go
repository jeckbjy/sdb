package sql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jeckbjy/sdb/engine/comm"

	"github.com/jeckbjy/sdb/engine"
)

func Register() {
	engine.Register("sql", New)
}

func New() engine.IClient {
	return &_Client{}
}

// TODO:如何支持json
type _Client struct {
	db    *sql.DB
	codec engine.ICodec
}

func (c *_Client) Name() string {
	return "sql"
}

func (c *_Client) Open(opts *engine.OpenOptions) error {
	db, err := sql.Open(opts.Driver, opts.URI)
	if err != nil {
		return err
	}

	if opts.Codec == nil {
		opts.Codec = comm.DefaultCodec()
	}
	c.db = db

	return nil
}

func (c *_Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}

	return nil
}

func (c *_Client) Ping() error {
	if c.db != nil {
		return c.db.Ping()
	}

	return nil
}

func (c *_Client) Database(name string) (engine.IDatabase, error) {
	if c.db != nil {
		return &_Database{db: c.db, codec: c.codec}, nil
	}

	return nil, errors.New("bad db")
}

func (c *_Client) Drop(name string) error {
	query := fmt.Sprintf("DROP database %s", name)
	_, err := c.db.Exec(query)
	return err
}
