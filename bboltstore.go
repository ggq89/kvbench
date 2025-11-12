package kvbench

import (
	"bytes"

	"go.etcd.io/bbolt"
)

var bboltBucket = []byte("keys")

type bboltStore struct {
	db *bbolt.DB
}

func NewBboltStore(path string, fsync bool) (Store, error) {
	if path == ":memory:" {
		return nil, ErrMemoryNotAllowed
	}
	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}
	db.NoSync = !fsync

	if err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bboltBucket)
		return err
	}); err != nil {
		db.Close()
		return nil, err
	}
	return &bboltStore{
		db: db,
	}, nil
}

func (s *bboltStore) Close() error {
	s.db.Close()
	return nil
}

func (s *bboltStore) PSet(keys, values [][]byte) error {
	return s.db.Batch(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bboltBucket)
		for i := 0; i < len(keys); i++ {
			if err := b.Put(keys[i], values[i]); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *bboltStore) PGet(keys [][]byte) ([][]byte, []bool, error) {
	var values [][]byte
	var oks []bool
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bboltBucket)
		for i := 0; i < len(keys); i++ {
			v := b.Get(keys[i])
			values = append(values, v)
			oks = append(oks, v != nil)
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return values, oks, nil
}

func (s *bboltStore) Set(key, value []byte) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(bboltBucket).Put(key, value)
	})
}

func (s *bboltStore) Get(key []byte) ([]byte, bool, error) {
	var v []byte
	err := s.db.View(func(tx *bbolt.Tx) error {
		v = tx.Bucket(bboltBucket).Get(key)
		return nil
	})
	return v, v != nil, err
}

func (s *bboltStore) Del(key []byte) (bool, error) {
	err := s.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(bboltBucket).Delete(key)
	})
	return err == nil, err
}

func (s *bboltStore) Keys(pattern []byte, limit int, withvalues bool) ([][]byte, [][]byte, error) {
	var keys [][]byte
	var vals [][]byte
	err := s.db.View(func(tx *bbolt.Tx) error {
		// Assume bucket exists and has keys
		c := tx.Bucket([]byte(bboltBucket)).Cursor()

		prefix := pattern
		for key, value := c.Seek(prefix); key != nil && bytes.HasPrefix(key, prefix); key, value = c.Next() {
			keys = append(keys, key)
			if withvalues {
				vals = append(vals, value)
			}
		}
		return nil
	})

	return keys, vals, err
}

func (s *bboltStore) FlushDB() error {
	return s.db.Sync()
}
