package kvbench

import (
	rosedb "github.com/rosedblabs/rosedb/v2"
)

type rosedbStore struct {
	db *rosedb.DB
}

func rosedbKey(key []byte) []byte {
	r := make([]byte, len(key)+1)
	r[0] = 'k'
	copy(r[1:], key)
	return r
}

func NewRosedbStore(path string, fsync bool) (Store, error) {
	if path == ":memory:" {
		return nil, ErrMemoryNotAllowed
	}

	options := rosedb.DefaultOptions
	options.DirPath = path
	options.Sync = fsync

	db, err := rosedb.Open(options)
	if err != nil {
		return nil, err
	}

	return &rosedbStore{
		db: db,
	}, nil
}

func (s *rosedbStore) Close() error {
	s.db.Close()
	return nil
}

func (s *rosedbStore) PSet(keys, vals [][]byte) error {
	batch := s.db.NewBatch(rosedb.DefaultBatchOptions)
	for i, k := range keys {
		batch.Put(k, vals[i])
	}

	return batch.Commit()
}

func (s *rosedbStore) PGet(keys [][]byte) ([][]byte, []bool, error) {
	var vals = make([][]byte, len(keys))
	var oks = make([]bool, len(keys))
	var err error

	opt := rosedb.DefaultBatchOptions
	opt.ReadOnly = true
	batch := s.db.NewBatch(opt)
	for i, k := range keys {
		vals[i], err = batch.Get(k)
		oks[i] = (err == nil)
	}

	return vals, oks, err
}

func (s *rosedbStore) Set(key, value []byte) error {
	return s.db.Put(key, value)
}

func (s *rosedbStore) Get(key []byte) ([]byte, bool, error) {
	v, err := s.db.Get(key)
	return v, v != nil, err
}

func (s *rosedbStore) Del(key []byte) (bool, error) {
	err := s.db.Delete(key)
	return err == nil, err
}

func (s *rosedbStore) Keys(pattern []byte, limit int, withvals bool) ([][]byte, [][]byte, error) {
	var keys [][]byte
	var vals [][]byte

	iterOpts := rosedb.DefaultIteratorOptions
	iterOpts.Prefix = pattern
	iter := s.db.NewIterator(iterOpts)
	defer iter.Close()

	for iter.Rewind(); iter.Valid(); iter.Next() {
		item := iter.Item()
		if item != nil {
			keys = append(keys, item.Key)
			if withvals {
				vals = append(vals, item.Value)
			}
		}
	}
	if err := iter.Err(); err != nil {
		return nil, nil, err
	}

	return keys, vals, nil
}

func (s *rosedbStore) FlushDB() error {
	return s.db.Sync()
}
