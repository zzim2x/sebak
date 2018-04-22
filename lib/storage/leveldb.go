package storage

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spikeekips/sebak/lib/util"
	"github.com/syndtr/goleveldb/leveldb"
	leveldbStorage "github.com/syndtr/goleveldb/leveldb/storage"
	leveldbUtil "github.com/syndtr/goleveldb/leveldb/util"
)

type LevelDBBackend struct {
	DB *leveldb.DB
}

func (st *LevelDBBackend) Init(config StorageConfig) (err error) {
	var sto leveldbStorage.Storage
	if path, ok := config["path"]; !ok {
		err = fmt.Errorf("`path` is not missing")
	} else if path == "<memory>" {
		sto = leveldbStorage.NewMemStorage()
	} else {
		if sto, err = leveldbStorage.OpenFile(path, false); err != nil {
			return
		}
	}

	if st.DB, err = leveldb.Open(sto, nil); err != nil {
		return
	}

	return
}

func (st *LevelDBBackend) Close() error {
	return st.DB.Close()
}

func (st *LevelDBBackend) Has(k string) (bool, error) {
	return st.DB.Has([]byte(k), nil)
}

func (st *LevelDBBackend) GetRaw(k string) (b []byte, err error) {
	var exists bool
	if exists, err = st.Has(k); !exists || err != nil {
		if !exists {
			err = fmt.Errorf("key, '%s' does not exists", k)
		}
		return
	}

	b, err = st.DB.Get([]byte(k), nil)

	return
}

func (st *LevelDBBackend) Get(k string, i interface{}) (err error) {
	var b []byte
	if b, err = st.GetRaw(k); err != nil {
		return
	}

	if err = json.Unmarshal(b, &i); err != nil {
		return
	}

	return
}

func (st *LevelDBBackend) New(k string, v interface{}) (err error) {
	var encoded []byte
	serializable, ok := v.(Serializable)
	if ok {
		encoded, err = serializable.Serialize()
	} else {
		encoded, err = util.EncodeJSONValue(v)
	}
	if err != nil {
		return
	}

	var exists bool
	if exists, err = st.Has(k); exists || err != nil {
		if exists {
			err = fmt.Errorf("key, '%s' already exists", k)
		}
		return
	}

	err = st.DB.Put([]byte(k), encoded, nil)

	return
}

func (st *LevelDBBackend) News(vs ...Item) (err error) {
	if len(vs) < 1 {
		err = errors.New("empty values")
		return
	}

	var exists bool
	for _, v := range vs {
		if exists, err = st.Has(v.Key); exists || err != nil {
			if exists {
				err = fmt.Errorf("found existing key, '%s'", v.Key)
			}
			return
		}
	}

	batch := new(leveldb.Batch)
	for _, v := range vs {
		var encoded []byte
		if encoded, err = util.EncodeJSONValue(v); err != nil {
			return
		}

		batch.Put([]byte(v.Key), encoded)
	}

	err = st.DB.Write(batch, nil)

	return
}

func (st *LevelDBBackend) Set(k string, v interface{}) (err error) {
	var encoded []byte
	if encoded, err = util.EncodeJSONValue(v); err != nil {
		return
	}

	var exists bool
	if exists, err = st.Has(k); !exists || err != nil {
		if !exists {
			err = fmt.Errorf("key, '%s' does not exists", k)
		}
		return
	}

	err = st.DB.Put([]byte(k), encoded, nil)

	return
}

func (st *LevelDBBackend) Sets(vs ...Item) (err error) {
	if len(vs) < 1 {
		err = errors.New("empty values")
		return
	}

	var exists bool
	for _, v := range vs {
		if exists, err = st.Has(v.Key); !exists || err != nil {
			if !exists {
				err = fmt.Errorf("not found key, '%s'", v.Key)
			}
			return
		}
	}

	batch := new(leveldb.Batch)
	for _, v := range vs {
		var encoded []byte
		if encoded, err = util.EncodeJSONValue(v); err != nil {
			return
		}

		batch.Put([]byte(v.Key), encoded)
	}

	err = st.DB.Write(batch, nil)

	return
}

func (st *LevelDBBackend) Remove(k string) (err error) {
	var exists bool
	if exists, err = st.Has(k); !exists || err != nil {
		if !exists {
			err = fmt.Errorf("key, '%s' does not exists", k)
		}
		return
	}

	err = st.DB.Delete([]byte(k), nil)

	return
}

func (st *LevelDBBackend) GetIterator(prefix string, reverse bool) (func() (IterItem, bool), func()) {
	var dbRange *leveldbUtil.Range
	if len(prefix) > 0 {
		dbRange = leveldbUtil.BytesPrefix([]byte(prefix))
	}

	iter := st.DB.NewIterator(dbRange, nil)

	var funcNext func() bool
	var hasUnsent bool
	if reverse {
		if !iter.Last() {
			iter.Release()
			return (func() (IterItem, bool) { return IterItem{}, false }), (func() {})
		}
		funcNext = iter.Prev
		hasUnsent = true
	} else {
		funcNext = iter.Next
		hasUnsent = false
	}

	var n int64 = 0
	return (func() (IterItem, bool) {
			if hasUnsent {
				hasUnsent = false
				return IterItem{N: n, Key: iter.Key(), Value: iter.Value()}, true
			}

			if !funcNext() {
				iter.Release()
				return IterItem{}, false
			}

			n += 1
			return IterItem{N: n, Key: iter.Key(), Value: iter.Value()}, true
		}),
		(func() {
			iter.Release()
		})
}
