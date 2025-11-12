package kvbench

import (
	lotusdb "github.com/lotusdblabs/lotusdb/v2"
)

type lotusdbStore struct {
	db *lotusdb.DB
}

func lotusdbKey(key []byte) []byte {
	r := make([]byte, len(key)+1)
	r[0] = 'k'
	copy(r[1:], key)
	return r
}

func NewLotusdbStore(path string, fsync bool) (Store, error) {
	if path == ":memory:" {
		return nil, ErrMemoryNotAllowed
	}

	opts := lotusdb.DefaultOptions
	opts.Sync = fsync
	opts.DirPath = path

	db, err := lotusdb.Open(opts)
	if err != nil {
		return nil, err
	}

	return &lotusdbStore{
		db: db,
	}, nil
}

func (s *lotusdbStore) Close() error {
	s.db.Close()
	return nil
}

func (s *lotusdbStore) PSet(keys, vals [][]byte) error {
	batch := s.db.NewBatch(lotusdb.DefaultBatchOptions)
	for i, k := range keys {
		batch.Put(k, vals[i])
	}

	return batch.Commit()
}

func (s *lotusdbStore) PGet(keys [][]byte) ([][]byte, []bool, error) {
	var vals = make([][]byte, len(keys))
	var oks = make([]bool, len(keys))
	var err error

	batch := s.db.NewBatch(lotusdb.DefaultBatchOptions)
	for i, k := range keys {
		vals[i], err = batch.Get(k)
		oks[i] = (err == nil)
	}

	return vals, oks, err
}

func (s *lotusdbStore) Set(key, value []byte) error {
	return s.db.Put(key, value)
}

func (s *lotusdbStore) Get(key []byte) ([]byte, bool, error) {
	v, err := s.db.Get(key)
	return v, v != nil, err
}

func (s *lotusdbStore) Del(key []byte) (bool, error) {
	err := s.db.Delete(key)
	return err == nil, err
}

func (s *lotusdbStore) Keys(pattern []byte, limit int, withvals bool) ([][]byte, [][]byte, error) {
	var keys [][]byte
	var vals [][]byte

	iter, err := s.db.NewIterator(lotusdb.IteratorOptions{Prefix: pattern, Reverse: false})
	if err != nil {
		return nil, nil, err
	}
	defer iter.Close()

	for iter.Valid() {
		keys = append(keys, iter.Key())
		if withvals {
			vals = append(vals, iter.Value())
		}
		iter.Next()
	}

	return keys, vals, nil
}

func (s *lotusdbStore) FlushDB() error {
	return s.db.Sync()
}
