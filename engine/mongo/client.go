package mongo

import (
	"context"
	"time"

	"github.com/jeckbjy/sdb/engine"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Register() {
	engine.Register("mongo", New)
}

func New() engine.IClient {
	return &_Client{}
}

type _Client struct {
	client *mongo.Client
}

func (c *_Client) Name() string {
	return "mongo"
}

func (c *_Client) Open(opts *engine.OpenOptions) error {
	// create from outside
	if opts.Client != nil {
		if client, ok := opts.Client.(*mongo.Client); ok {
			c.client = client
			return nil
		}
	}

	// 默认连接本地
	if opts.URI == engine.DefaultURI {
		opts.URI = "mongodb://localhost:27017"
	}

	// 其他一些配置,比如max pool size, read ref等
	// mongo://localhost:27017
	o := options.Client().ApplyURI(opts.URI)

	client, err := mongo.NewClient(o)
	if err != nil {
		return err
	}

	c.client = client
	if err := c.Ping(); err != nil {
		return err
	}

	return nil
}

func (c *_Client) Close() error {
	if c.client != nil {
		return c.client.Disconnect(nil)
	}

	return nil
}

func (c *_Client) Ping() error {
	ctx, cb := context.WithTimeout(context.Background(), time.Second*10)
	defer cb()

	return c.client.Ping(ctx, nil)
}

func (c *_Client) Database(name string) (engine.IDatabase, error) {
	db := c.client.Database(name)
	return &_Database{db: db}, nil
}

func (c *_Client) Drop(name string) error {
	return c.client.Database(name).Drop(nil)
}
