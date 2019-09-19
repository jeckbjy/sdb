package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/jeckbjy/sdb/engine/comm"

	"github.com/jeckbjy/sdb/engine"
)

type _Database struct {
	db    *sql.DB
	codec engine.ICodec
}

func (d *_Database) EnsureIndex(bucket string, keys interface{}, opts *engine.IndexOptions) (string, error) {
	var name, cols string
	switch keys.(type) {
	case string:
		name = keys.(string)
		cols = name
	case []string:
		k := keys.([]string)
		name = strings.Join(k, "_")
		cols = strings.Join(k, ",")
	default:
		return "", errors.New("ensure index fail,not support")
	}

	if opts.Name != "" {
		name = opts.Name
	}

	var indexType string
	if opts.Unique {
		indexType = "UNIQUE INDEX"
	} else {
		indexType = "INDEX"
	}

	// mysql not support IF NOT EXISTS ??
	query := fmt.Sprintf("CREATE %s IF NOT EXISTS %s ON %s (%s)", indexType, name, bucket, cols)
	_, err := d.db.Exec(query)
	return name, err
}

func (d *_Database) DropIndex(bucket string, name string) error {
	_, err := d.db.Exec("DROP INDEX IF EXISTS %s ON %s", name, bucket)
	return err
}

func (d *_Database) Indexes(bucket string) ([]engine.Index, error) {
	// MYSQL: SHOW INDEXES FROM table_name IN database_name;
	// PG: SELECT tablename,indexname,indexdef FROM pg_indexes WHERE schemaname = 'public' ORDER BY tablename, indexname;
	panic("implement me")
}

func (d *_Database) Create(bucket string, doc interface{}) error {
	// migrate
	panic("implement me")
}

func (d *_Database) Drop(bucket string) error {
	_, err := d.db.Exec("DROP TABLE %s", bucket)
	return err
}

func (d *_Database) Insert(bucket string, doc interface{}, opts *engine.InsertOptions) (*engine.InsertResult, error) {
	// INSERT INTO table_name (col1, col2,...) VALUES (val1, val2,...)
	res := &engine.InsertResult{}
	rows := make([]map[string]interface{}, 0)
	if opts.One {
		rows = append(rows, comm.ToMap(d.codec, doc))
	} else {
		slice := reflect.ValueOf(doc)
		for i := 0; i < slice.Len(); i++ {
			item := slice.Index(i)
			rows = append(rows, comm.ToMap(d.codec, item.Interface()))
		}
	}

	for _, row := range rows {
		id, ok := row["_id"]
		if !ok {
			id = comm.NewID()
			row["_id"] = id
		}

		// 如何支持json类型呢?
		cb := builder{}
		vb := builder{}
		for k, v := range row {
			cb.WriteBy(k, ',')
			vb.WriteBy(toValue(v), ',')
		}

		_, err := d.db.Exec("INSERT INTO %s (%s) VALUES (%s)", bucket, cb.String(), vb.String())
		if err != nil {
			return nil, err
		}

		res.InsertedIDs = append(res.InsertedIDs, id)
	}

	return res, nil
}

func (d *_Database) Delete(bucket string, filter engine.ICond, opts *engine.DeleteOptions) (*engine.DeleteResult, error) {
	// delete from %s where %s limit
	where, err := toWhere(filter)
	if err != nil {
		return nil, err
	}

	b := builder{}
	b.Writef("DELETE FROM %s", bucket)
	b.Write(where)
	if opts.One {
		b.Write("LIMIT 1")
	}

	res, err := d.db.Exec(b.String())
	if err != nil {
		return nil, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	return &engine.DeleteResult{DeletedCount: count}, nil
}

func (d *_Database) Update(bucket string, filter engine.ICond, update interface{}, opts *engine.UpdateOptions) (*engine.UpdateResult, error) {
	// UPDATE table_name SET 列名称 = 新值 WHERE 列名称 = 某值 [LIMIT 1]
	where, err := toWhere(filter)
	if err != nil {
		return nil, err
	}

	m := comm.ToMap(d.codec, update)
	b := builder{}
	for k, v := range m {
		b.WriteBy(fmt.Sprintf("%s=%s", k, toValue(v)), ',')
	}

	q := builder{}
	q.Writef("UPDATE %s SET VALUES (%s)", bucket, b.String())
	q.Write(where)
	if opts.One {
		q.Write("LIMIT 1")
	}

	if r, err := d.db.Exec(q.String()); err != nil {
		return nil, err
	} else {
		count, _ := r.RowsAffected()
		res := &engine.UpdateResult{}
		res.MatchedCount = count
		return res, nil
	}
}

func (d *_Database) Query(result interface{}, bucket string, filter engine.ICond, opts *engine.QueryOptions) error {
	where, err := toWhere(filter)
	if err != nil {
		return err
	}

	b := builder{}
	b.Writef("SELECT %s FROM %s")
	b.Write(where)
	if opts.One {
		b.Write("LIMIT 1")
	}

	if opts.One {
		row := d.db.QueryRow(b.String())
		return row.Scan(result)
	} else {
		rows, err := d.db.Query(b.String())
		if err != nil {
			return err
		}

		resultv := reflect.ValueOf(result)
		slicev := resultv.Elem()
		elemt := slicev.Type().Elem()
		for rows.Next() {
			elemp := reflect.New(elemt)
			if err := rows.Scan(elemp.Interface()); err != nil {
				return err
			}
			slicev = reflect.Append(slicev, elemp.Elem())
		}

		resultv.Elem().Set(slicev)
	}

	return nil
}
