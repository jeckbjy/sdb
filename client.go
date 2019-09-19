package sdb

import (
	"errors"

	"github.com/jeckbjy/sdb/engine"
)

var ErrDriverNotFound = errors.New("driver not found")

const (
	MONGO = "mongo"
	BOLT  = "bolt"
	SQL   = "sql"
)

func New(driver string, opts ...OpenOption) (*Client, error) {
	eng := engine.New(driver)
	if eng == nil {
		return nil, ErrDriverNotFound
	}

	c := &Client{client: eng}
	e := c.Open(opts...)
	return c, e
}

type Client struct {
	client engine.IClient
}

func (c *Client) Open(opts ...OpenOption) error {
	o := engine.OpenOptions{}
	o.URI = DefaultURI
	for _, fn := range opts {
		fn(&o)
	}

	return c.client.Open(&o)
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Ping() error {
	return c.client.Ping()
}

func (c *Client) Database(name string) (*Database, error) {
	db, err := c.client.Database(name)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (c *Client) Drop(name string) error {
	return c.client.Drop(name)
}
