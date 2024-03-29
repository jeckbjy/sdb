# sdb
simple database client wrapper for mongo, sql, or kv store.

封装常用的数据库操作Insert,Delete,Update,Query,Index操作,以保障不同的底层引擎对相同的操作会有一致的行为  
目前主要封装了三种数据库,mongo,sql,bolt，bolt并没有索引，完全是暴力全遍历查询，只能用作本地测试使用，API设计上主要向mongo靠近

目前仅仅是粗略的实现了一下，还没有细致的测试

用法:
```go

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

// 分页查询
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
```