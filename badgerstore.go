package kvbench

import (
	"github.com/dgraph-io/badger/v4"
)

type badgerStore struct {
	db *badger.DB
}

func NewBadgerStore(path string, fsync bool) (Store, error) {
	opts := badger.DefaultOptions(path)
	opts = opts.WithLoggingLevel(badger.ERROR)
	if path == ":memory:" {
		opts.InMemory = true
	}
	opts.Logger = nil

	opts.SyncWrites = fsync
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &badgerStore{
		db: db,
	}, nil
}

func (s *badgerStore) Close() error {
	s.db.Close()
	return nil
}

func (s *badgerStore) PSet(keys, vals [][]byte) error {
	wb := s.db.NewWriteBatch()
	for i := range keys {
		err := wb.Set(keys[i], vals[i])
		if err != nil {
			return err
		}
	}
	return wb.Flush()
}

func (s *badgerStore) PGet(keys [][]byte) ([][]byte, []bool, error) {
	var vals = make([][]byte, len(keys))
	var oks = make([]bool, len(keys))

	err := s.db.View(func(txn *badger.Txn) error {
		for i, k := range keys {
			item, err := txn.Get(k)
			if err == nil {
				v, err := item.ValueCopy(nil)
				if err == nil {
					vals[i] = v
					oks[i] = true
				}
			}
		}
		return nil
	})

	return vals, oks, err
}

func (s *badgerStore) Set(key, value []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (s *badgerStore) Get(key []byte) ([]byte, bool, error) {
	var v []byte

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err == nil {
			err = item.Value(func(value []byte) error {
				v = value
				return nil
			})
			// v, err = item.ValueCopy(nil)
		}
		return err
	})

	return v, v != nil, err
}

func (s *badgerStore) Del(key []byte) (bool, error) {
	err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
	return err == nil, err
}

func (s *badgerStore) Keys(pattern []byte, limit int, withvals bool) ([][]byte, [][]byte, error) {
	var keys [][]byte
	var vals [][]byte

	err := s.db.View(func(txn *badger.Txn) error {
		opt := badger.DefaultIteratorOptions
		opt.Prefix = pattern
		it := txn.NewIterator(opt)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			keys = append(keys, item.Key())
			if withvals {
				item.Value(func(v []byte) error {
					vals = append(vals, v)
					return nil
				})
			}

		}

		return nil
	})

	return keys, vals, err
}

func (s *badgerStore) FlushDB() error {
	return s.db.Sync()
}
