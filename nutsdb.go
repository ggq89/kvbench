package kvbench

import (
	"github.com/nutsdb/nutsdb"
)

const nutsdbBucket = "keys"

type nutsdbStore struct {
	// mu sync.RWMutex
	db *nutsdb.DB
}

const defaultSegmentSize int64 = 256 * nutsdb.MB

// DefaultOptions represents the default options.
var pptions = func() nutsdb.Options {
	opt := nutsdb.DefaultOptions
	opt.EntryIdxMode = nutsdb.HintKeyAndRAMIdxMode
	opt.SegmentSize = defaultSegmentSize
	opt.NodeNum = 1
	opt.RWMode = nutsdb.FileIO
	opt.SyncEnable = false
	return opt
}()

// func nutsdbKey(key []byte) []byte {
// 	r := make([]byte, len(key)+1)
// 	r[0] = 'k'
// 	copy(r[1:], key)
// 	return r
// }

func NewNutsdbStore(path string, fsync bool) (Store, error) {
	if path == ":memory:" {
		return nil, ErrMemoryNotAllowed
	}

	opt := pptions
	opt.SyncEnable = fsync
	opt.Dir = path

	db, err := nutsdb.Open(opt)
	if err != nil {
		return nil, err
	}

	// err = db.Update(func(tx *nutsdb.Tx) error {
	// 	// you should call Bucket with data structure and the name of bucket first
	// 	return tx.NewKVBucket(nutsdbBucket)
	// })
	// if err != nil {
	// 	return nil, err
	// }

	return &nutsdbStore{
		db: db,
	}, nil
}

func (s *nutsdbStore) Close() error {
	s.db.Close()
	return nil
}

func (s *nutsdbStore) PSet(keys, vals [][]byte) error {
	return s.db.Update(func(tx *nutsdb.Tx) error {
		for i, k := range keys {
			tx.Put(nutsdbBucket, k, vals[i], nutsdb.Persistent)
		}

		return nil
	})
}

func (s *nutsdbStore) PGet(keys [][]byte) ([][]byte, []bool, error) {
	var vals = make([][]byte, len(keys))
	var oks = make([]bool, len(keys))

	var err error

	s.db.View(func(tx *nutsdb.Tx) error {
		for i, k := range keys {
			e, err := tx.Get(nutsdbBucket, k)
			if e != nil {
				vals[i] = e.Value
			}

			oks[i] = (err == nil)
		}

		return nil
	})

	return vals, oks, err
}

func (s *nutsdbStore) Set(key, value []byte) error {
	err := s.db.Update(func(tx *nutsdb.Tx) error {
		e := tx.Put(nutsdbBucket, key, value, nutsdb.Persistent)
		return e
	})

	return err
}

func (s *nutsdbStore) Get(key []byte) ([]byte, bool, error) {
	var v []byte
	var ok bool
	var err error

	s.db.View(func(tx *nutsdb.Tx) error {
		e, err := tx.Get(nutsdbBucket, key)
		if e != nil {
			v = e.Value
		}
		ok = err == nil
		return err
	})

	return v, ok, err
}

func (s *nutsdbStore) Del(key []byte) (bool, error) {
	err := s.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Delete(nutsdbBucket, key)
	})

	return err == nil, err
}

func (s *nutsdbStore) Keys(pattern []byte, limit int, withvals bool) ([][]byte, [][]byte, error) {
	var keys [][]byte
	var vals [][]byte

	err := s.db.View(func(tx *nutsdb.Tx) error {
		entries, err := tx.PrefixScan(nutsdbBucket, pattern, 0, nutsdb.ScanNoLimit)
		if err != nil {
			return err
		}

		for i, entry := range entries {
			keys[i] = entry.Key
			vals[i] = entry.Value
		}

		// iterator := nutsdb.NewIterator(tx, nutsdbBucket, nutsdb.IteratorOptions{Reverse: false})

		// for iterator.Valid() {
		// 	keys = append(keys, bcopy(iterator.Key()))
		// 	if withvals {
		// 		value, err := iterator.Value()
		// 		if err != nil {
		// 			return err
		// 		}
		// 		vals = append(vals, bcopy(value))
		// 	}

		// 	if !iterator.Next() {
		// 		break
		// 	}
		// }

		return nil
	})

	return keys, vals, err
}

func (s *nutsdbStore) FlushDB() error {
	return s.db.Close()
}
