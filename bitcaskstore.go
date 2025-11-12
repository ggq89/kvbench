package kvbench

import (
	bitcask "go.mills.io/bitcask/v2"
)

type bitcaskStore struct {
	db *bitcask.Bitcask
}

func bitcaskKey(key []byte) []byte {
	r := make([]byte, len(key)+1)
	r[0] = 'k'
	copy(r[1:], key)
	return r
}

func NewBitcaskStore(path string, fsync bool) (Store, error) {
	if path == ":memory:" {
		return nil, ErrMemoryNotAllowed
	}

	db, err := bitcask.Open(path, bitcask.WithSyncWrites(fsync))
	if err != nil {
		return nil, err
	}

	return &bitcaskStore{
		db: db,
	}, nil
}

func (s *bitcaskStore) Close() error {
	s.db.Close()
	return nil
}

func (s *bitcaskStore) PSet(keys, vals [][]byte) error {
	batch := s.db.Batch()
	for i, k := range keys {
		batch.Put(k, vals[i])
	}

	return s.db.WriteBatch(batch)
}

func (s *bitcaskStore) PGet(keys [][]byte) ([][]byte, []bool, error) {
	var vals = make([][]byte, len(keys))
	var oks = make([]bool, len(keys))
	var err error

	for i, k := range keys {
		vals[i], err = s.db.Get(k)
		oks[i] = (err == nil)
	}

	return vals, oks, err
}

func (s *bitcaskStore) Set(key, value []byte) error {
	return s.db.Put(key, value)
}

func (s *bitcaskStore) Get(key []byte) ([]byte, bool, error) {
	v, err := s.db.Get(key)
	return v, v != nil, err
}

func (s *bitcaskStore) Del(key []byte) (bool, error) {
	err := s.db.Delete(key)
	return err == nil, err
}

func (s *bitcaskStore) Keys(pattern []byte, limit int, withvals bool) ([][]byte, [][]byte, error) {
	var keys [][]byte
	var vals [][]byte

	s.db.Scan(pattern, func(k bitcask.Key) error {
		keys = append(keys, k)
		if withvals {
			v, err := s.db.Get(k)
			if err != nil {
				return err
			}
			vals = append(vals, v)
		}
		return nil
	})

	return keys, vals, nil
}

func (s *bitcaskStore) FlushDB() error {
	return s.db.Sync()
}
