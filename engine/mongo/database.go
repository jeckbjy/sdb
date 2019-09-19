package mongo

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/jeckbjy/sdb/engine"
	"go.mongodb.org/mongo-driver/mongo"
)

// https://godoc.org/github.com/globalsign/mgo
type _Database struct {
	db *mongo.Database
}

func (d *_Database) EnsureIndex(bucket string, keys interface{}, opts *engine.IndexOptions) (string, error) {
	o := options.Index()
	if opts.Name != "" {
		o.SetName(opts.Name)
	}

	if opts.Background {
		o.SetBackground(true)
	}

	if opts.Sparse {
		o.SetSparse(true)
	}

	if opts.Unique {
		o.SetUnique(true)
	}

	model := mongo.IndexModel{}
	model.Options = o

	switch keys.(type) {
	case string:
		model.Keys = bson.M{keys.(string): 1}
	case []string:
		kk := bson.M{}
		for _, k := range keys.([]string) {
			kk[k] = 1
		}
		model.Keys = kk
	default:
		model.Keys = keys
	}

	indexes := d.db.Collection(bucket).Indexes()
	return indexes.CreateOne(context.Background(), model)
}

func (d *_Database) DropIndex(bucket string, name string) error {
	//o := options.DropIndexes()
	indexes := d.db.Collection(bucket).Indexes()
	_, err := indexes.DropOne(context.Background(), name)
	return err
}

func (d *_Database) Indexes(bucket string) ([]engine.Index, error) {
	indexes := d.db.Collection(bucket).Indexes()
	cursor, err := indexes.List(context.Background())
	if err != nil {
		return nil, err
	}

	results := make([]engine.Index, 0)
	for cursor.Next(nil) {
		index := engine.Index{}
		if err := cursor.Decode(index); err != nil {
			return nil, err
		}
	}

	return results, nil
}

func (d *_Database) Create(bucket string, doc interface{}) error {
	return nil
}

func (d *_Database) Drop(bucket string) error {
	return d.db.Collection(bucket).Drop(nil)
}

func (d *_Database) Insert(bucket string, doc interface{}, opts *engine.InsertOptions) (*engine.InsertResult, error) {
	coll := d.db.Collection(bucket)
	if opts.One {
		res, err := coll.InsertOne(nil, doc)
		if err != nil {
			return nil, err
		}
		return &engine.InsertResult{InsertedIDs: []interface{}{res.InsertedID}}, nil
	} else {
		res, err := coll.InsertMany(nil, doc.([]interface{}))
		if err != nil {
			return nil, err
		}

		return &engine.InsertResult{InsertedIDs: res.InsertedIDs}, nil
	}
}

func (d *_Database) Delete(bucket string, filter engine.ICond, opts *engine.DeleteOptions) (*engine.DeleteResult, error) {
	var res *mongo.DeleteResult
	var err error

	coll := d.db.Collection(bucket)
	query, err := toBson(filter)
	if err != nil {
		return nil, err
	}

	if opts.One {
		res, err = coll.DeleteOne(nil, query)
	} else {
		res, err = coll.DeleteMany(nil, query)
	}
	if err != nil {
		return nil, err
	}

	return &engine.DeleteResult{DeletedCount: res.DeletedCount}, nil
}

func (d *_Database) Update(bucket string, filter engine.ICond, update interface{}, opts *engine.UpdateOptions) (*engine.UpdateResult, error) {
	var res *mongo.UpdateResult
	var err error

	coll := d.db.Collection(bucket)
	query, err := toBson(filter)
	if err != nil {
		return nil, err
	}

	o := options.Update()
	if opts.Upsert {
		o.SetUpsert(opts.Upsert)
	}

	if opts.One {
		res, err = coll.UpdateOne(nil, bson.M{"$set": query}, update, o)
	} else {
		res, err = coll.UpdateMany(nil, bson.M{"$set": query}, update, o)
	}
	if err != nil {
		return nil, err
	}

	return &engine.UpdateResult{MatchedCount: res.MatchedCount, ModifiedCount: res.ModifiedCount, UpsertedCount: res.UpsertedCount, UpsertedID: res.UpsertedID}, nil
}

func (d *_Database) Query(result interface{}, bucket string, filter engine.ICond, opts *engine.QueryOptions) error {
	coll := d.db.Collection(bucket)
	query, err := toBson(filter)
	if err != nil {
		return err
	}
	if opts.One {
		o := options.FindOne()
		if opts.Skip > 0 {
			o.SetSkip(opts.Skip)
		}
		if opts.Sort != nil {
			o.SetSort(opts.Sort)
		}
		if opts.Projection != nil {
			o.SetProjection(opts.Projection)
		}

		res := coll.FindOne(nil, query, o)
		if res.Err() != nil {
			return res.Err()
		}

		return res.Decode(result)
	} else {
		o := options.Find()
		if opts.Skip > 0 {
			o.SetSkip(opts.Skip)
		}
		if opts.Limit > 0 {
			o.SetLimit(opts.Limit)
		}
		if opts.Sort != nil {
			o.SetSort(opts.Sort)
		}
		if opts.Projection != nil {
			o.SetProjection(opts.Projection)
		}

		cursor, err := coll.Find(nil, query, o)
		if err != nil {
			return err
		}

		resultv := reflect.ValueOf(result)
		slicev := resultv.Elem()
		elemt := slicev.Type().Elem()
		for cursor.Next(nil) {
			elemp := reflect.New(elemt)
			if err := cursor.Decode(elemp.Interface()); err != nil {
				return err
			}

			slicev = reflect.Append(slicev, elemp.Elem())
		}

		resultv.Elem().Set(slicev)

		return nil
	}
}
