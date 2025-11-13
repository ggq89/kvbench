package kvbench

import (
	"bytes"

	"github.com/Data-Corruption/lmdb-go/lmdb"
	"github.com/Data-Corruption/lmdb-go/lmdbscan"
)

const lmdbBucket = "keys"

type lmdbStore struct {
	env *lmdb.Env
	dbi lmdb.DBI
}

func NewLmdbStore(path string, fsync bool) (Store, error) {
	if path == ":memory:" {
		return nil, ErrMemoryNotAllowed
	}

	env, err := lmdb.NewEnv()
	if err != nil {
		return nil, err
	}
	if !fsync {
		env.SetFlags(lmdb.NoSync)
	}

	err = env.Open(path, 0, 0644)
	if err != nil {
		return nil, err
	}

	var dbi lmdb.DBI
	err = env.Update(func(txn *lmdb.Txn) (err error) {
		dbi, err = txn.OpenDBI(lmdbBucket, lmdb.Create)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &lmdbStore{
		env: env,
		dbi: dbi,
	}, nil
}

func (s *lmdbStore) Close() error {
	s.env.CloseDBI(s.dbi)
	return s.env.Close()
}

func (s *lmdbStore) PSet(keys, vals [][]byte) error {
	err := s.env.Update(func(txn *lmdb.Txn) (err error) {
		for i, k := range keys {
			err = txn.Put(s.dbi, k, vals[i], 0)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *lmdbStore) PGet(keys [][]byte) ([][]byte, []bool, error) {
	var vals = make([][]byte, len(keys))
	var oks = make([]bool, len(keys))

	err := s.env.View(func(txn *lmdb.Txn) (err error) {
		for i, k := range keys {
			e, err := txn.Get(s.dbi, k)
			if e != nil {
				vals[i] = e
			}
			oks[i] = (err == nil)
		}
		return nil
	})

	return vals, oks, err
}

func (s *lmdbStore) Set(key, value []byte) error {
	err := s.env.Update(func(txn *lmdb.Txn) (err error) {
		err = txn.Put(s.dbi, key, value, 0)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *lmdbStore) Get(key []byte) ([]byte, bool, error) {
	var v []byte
	var ok bool
	var err error

	err = s.env.View(func(txn *lmdb.Txn) (err error) {
		e, err := txn.Get(s.dbi, key)
		if e != nil {
			v = e
		}
		ok = err == nil
		return err
	})

	return v, ok, err
}

func (s *lmdbStore) Del(key []byte) (bool, error) {
	err := s.env.Update(func(txn *lmdb.Txn) (err error) {
		err = txn.Del(s.dbi, key, nil)
		if err != nil {
			return err
		}
		return nil
	})

	return err == nil, err
}

func (s *lmdbStore) Keys(pattern []byte, limit int, withvals bool) ([][]byte, [][]byte, error) {
	var keys [][]byte
	var vals [][]byte

	s.env.Close()

	err := s.env.View(func(txn *lmdb.Txn) (err error) {
		scanner := lmdbscan.New(txn, s.dbi)
		defer scanner.Close()

		scanner.Set(pattern, nil, lmdb.SetRange)
		for scanner.Scan() {
			if !bytes.HasPrefix(scanner.Key(), pattern) {
				break
			}
			keys = append(keys, scanner.Key())
			if withvals {
				vals = append(vals, scanner.Val())
			}
		}
		return scanner.Err()
	})
	if err != nil {
		return nil, nil, err
	}

	return keys, vals, err
}

func (s *lmdbStore) FlushDB() error {
	return s.env.Sync(true)
}
