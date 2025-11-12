package kvbench

import (
	"fmt"
	"io"

	"github.com/cockroachdb/pebble/v2"
)

type pebbleStore struct {
	db *pebble.DB
	wo *pebble.WriteOptions
}

func NewPebbleStore(path string, fsync bool) (Store, error) {
	if path == ":memory:" {
		return nil, ErrMemoryNotAllowed
	}

	opts := pebble.DefaultOptions()
	if !fsync {
		opts.DisableWAL = true
	}

	wo := &pebble.WriteOptions{}
	wo.Sync = fsync

	db, err := pebble.Open(path, opts)
	if err != nil {
		return nil, err
	}

	return &pebbleStore{
		db: db,
		wo: wo,
	}, nil
}

func (s *pebbleStore) Close() error {
	return s.db.Close()
}

func (s *pebbleStore) PSet(keys, vals [][]byte) error {
	wb := s.db.NewBatch()

	for i, k := range keys {
		wb.Set(k, vals[i], s.wo)
	}
	return wb.Commit(s.wo)
}

func (s *pebbleStore) PGet(keys [][]byte) ([][]byte, []bool, error) {
	var vals = make([][]byte, len(keys))
	var oks = make([]bool, len(keys))

	var err error
	var closer io.Closer
	for i, k := range keys {
		vals[i], closer, err = s.db.Get(k)
		oks[i] = (err == nil)
		if closer != nil {
			closer.Close()
		}
	}
	return vals, oks, err
}

func (s *pebbleStore) Set(key, value []byte) error {
	return s.db.Set(key, value, s.wo)
}

func (s *pebbleStore) Get(key []byte) ([]byte, bool, error) {
	v, closer, err := s.db.Get(key)
	if err != nil {
		if err != pebble.ErrNotFound {
			fmt.Println(err)
		}
		return nil, false, err
	}
	if closer != nil {
		closer.Close()
	}
	return v, v != nil, err
}

func (s *pebbleStore) Del(key []byte) (bool, error) {
	err := s.db.Delete(key, s.wo)
	return err == nil, err
}

func (s *pebbleStore) Keys(pattern []byte, limit int, withvals bool) ([][]byte, [][]byte, error) {
	var keys [][]byte
	var vals [][]byte

	keyUpperBound := func(b []byte) []byte {
		end := make([]byte, len(b))
		copy(end, b)
		for i := len(end) - 1; i >= 0; i-- {
			end[i] = end[i] + 1
			if end[i] != 0 {
				return end[:i+1]
			}
		}
		return nil // no upper-bound
	}

	prefixIterOptions := func(prefix []byte) *pebble.IterOptions {
		return &pebble.IterOptions{
			LowerBound: prefix,
			UpperBound: keyUpperBound(prefix),
		}
	}

	// 创建迭代器，指定范围为前缀
	iter, err := s.db.NewIter(prefixIterOptions(pattern))
	defer iter.Close()
	if err != nil {
		return nil, nil, err
	}

	// 遍历具有前缀的所有键值对
	for iter.First(); iter.Valid(); iter.Next() {
		key := iter.Key()
		keys = append(keys, key)
		if withvals {
			value := iter.Value()
			vals = append(vals, value)
		}
	}

	return keys, vals, nil
}

func (s *pebbleStore) FlushDB() error {
	return s.db.Flush()
}
