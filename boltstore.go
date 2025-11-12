package kvbench

import (
	"bytes"

	"github.com/boltdb/bolt"
)

var boltBucket = []byte("keys")

type boltStore struct {
	db *bolt.DB
}

func NewBoltStore(path string, fsync bool) (Store, error) {
	if path == ":memory:" {
		return nil, ErrMemoryNotAllowed
	}
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}
	db.NoSync = !fsync

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(boltBucket)
		return err
	}); err != nil {
		db.Close()
		return nil, err
	}

	return &boltStore{
		db: db,
	}, nil
}

func (s *boltStore) Close() error {
	s.db.Close()
	return nil
}

func (s *boltStore) PSet(keys, values [][]byte) error {
	return s.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(boltBucket)
		for i := 0; i < len(keys); i++ {
			if err := b.Put(keys[i], values[i]); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *boltStore) PGet(keys [][]byte) ([][]byte, []bool, error) {
	var values [][]byte
	var oks []bool
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(boltBucket)
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

func (s *boltStore) Set(key, value []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(boltBucket).Put(key, value)
	})
}

func (s *boltStore) Get(key []byte) ([]byte, bool, error) {
	var v []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		v = tx.Bucket(boltBucket).Get(key)
		return nil
	})
	return v, v != nil, err
}

func (s *boltStore) Del(key []byte) (bool, error) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(boltBucket).Delete(key)
	})
	return err == nil, err
}

func (s *boltStore) Keys(pattern []byte, limit int, withvalues bool) ([][]byte, [][]byte, error) {
	var keys [][]byte
	var vals [][]byte
	err := s.db.View(func(tx *bolt.Tx) error {
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

func (s *boltStore) FlushDB() error {
	return s.db.Sync()
}
