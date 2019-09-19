package tests

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/jeckbjy/sdb"
	"github.com/jeckbjy/sdb/engine/bolt"
)

type Foo struct {
	ID        string    `bson:"_id" json:"_id"`
	Name      string    `bson:"name" json:"name"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

func TestClient(t *testing.T) {
	bolt.Register()
	c, err := sdb.New("")
	if err != nil {
		t.Fatal(err)
	}

	db, err := c.Database("test")
	if err != nil {
		t.Fatal(err)
	}

	bucket := "foo"
	if _, err := db.Insert(bucket, &Foo{ID: "a", Name: "a", CreatedAt: time.Now()}); err != nil {
		t.Fatal(err)
	}

	result := Foo{}
	if err := db.Query(&result, bucket, sdb.Eq("_id", "a")); err != nil {
		t.Fatal(err)
	}

	t.Log(result)

	if _, err := db.Delete(bucket, sdb.Eq("_id", "a")); err != nil {
		t.Fatal(err)
	}

	if err := c.Drop("test"); err != nil {
		t.Fatal(err)
	}
}

func TestQuery(t *testing.T) {
	bolt.Register()
	c, err := sdb.New("")
	if err != nil {
		t.Fatal(err)
	}

	_ = c.Drop("test")

	db, err := c.Database("test")
	if err != nil {
		t.Fatal(err)
	}

	bucket := "foo"

	for i := 0; i < 10; i++ {
		_, err := db.Insert(bucket, &Foo{ID: fmt.Sprintf("%+v", i), Name: fmt.Sprintf("%+v", i), CreatedAt: time.Now().Add(time.Duration(i) * 10 * time.Minute)})
		if err != nil {
			t.Fatal(err)
		}
	}

	// page
	page := make([]Foo, 0)
	if err := db.Query(&page, bucket, sdb.Lt("_id", "4"), sdb.WithLimit(2), sdb.WithSkip(2), sdb.WithSort("createdAt", false)); err != nil {
		t.Fatal(err)
	}
	t.Log(page)
}

func TestJson(t *testing.T) {
	f := Foo{ID: "_id", Name: "name", CreatedAt: time.Now()}
	d, _ := json.Marshal(f)
	j := make(map[string]interface{})
	_ = json.Unmarshal(d, &j)
	t.Log(j)
}
