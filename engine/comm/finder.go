package comm

import (
	"reflect"
	"sort"

	"github.com/jeckbjy/sdb/engine"
)

type item struct {
	doc  map[string]interface{}
	data []byte
}

type Finder struct {
	items []item
}

func (f *Finder) Find(result interface{}, codec engine.ICodec, skip int64, limit int64, projection map[string]int) error {
	if skip > 0 {
		f.items = f.items[skip:]
	}

	if limit > 0 && limit < int64(len(f.items)) {
		f.items = f.items[:limit]
	}

	// TODO:
	if projection != nil {
	}

	resultv := reflect.ValueOf(result)
	slicev := resultv.Elem()
	elemt := slicev.Type().Elem()
	for _, v := range f.items {
		elemp := reflect.New(elemt)
		if err := codec.Decode(v.data, elemp.Interface()); err != nil {
			return err
		}

		slicev = reflect.Append(slicev, elemp.Elem())
	}

	resultv.Elem().Set(slicev)

	return nil
}

func (f *Finder) Push(doc map[string]interface{}, data []byte) {
	f.items = append(f.items, item{doc: doc, data: data})
}

func (f *Finder) Sort(sorts map[string]int) {
	if sorts == nil {
		return
	}
	//
	sort.Slice(f.items, func(i int, j int) bool {
		x := f.items[i]
		y := f.items[j]

		for k, v := range sorts {
			if v == 1 { // ascending
				if cmpLess(x.doc[k], y.doc[k]) {
					return true
				}
			} else { // descending -1
				if cmpGreat(x.doc[k], y.doc[k]) {
					return true
				}
			}
		}

		return false
	})
}
