package comm

import (
	"encoding/json"

	"github.com/jeckbjy/sdb/engine"
)

func DefaultCodec() engine.ICodec {
	return &JsonCodec{}
}

// 注:json编码会导致time.Time信息丢失,变为string
type JsonCodec struct {
}

func (c *JsonCodec) Name() string {
	return "json"
}

func (c *JsonCodec) Encode(doc interface{}) ([]byte, error) {
	return json.Marshal(doc)
}

func (c *JsonCodec) Decode(data []byte, doc interface{}) error {
	return json.Unmarshal(data, doc)
}
