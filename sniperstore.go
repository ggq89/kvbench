package kvbench

import (
	"sync"

	"github.com/recoilme/sniper"
)

type sniperStore struct {
	mu   sync.RWMutex
	path string
	db   *sniper.Store
}

func sniperKey(key []byte) []byte {
	r := make([]byte, len(key)+1)
	r[0] = 'k'
	copy(r[1:], key)
	return r
}

func NewSniperStore(path string, fsync bool) (Store, error) {
	if path == ":memory:" {
		return nil, ErrMemoryNotAllowed
	}
	db, err := sniper.Open(sniper.Dir(path))
	if err != nil {
		return nil, err
	}

	return &sniperStore{
		db:   db,
		path: path,
	}, nil
}

func (s *sniperStore) Close() error {
	s.db.Close()
	return nil
}

func (s *sniperStore) PSet(keys, values [][]byte) error {
	for i := 0; i < len(keys); i++ {
		if err := s.db.Set(sniperKey(keys[i]), values[i], 0); err != nil {
			return err
		}
	}
	return nil

}

func (s *sniperStore) PGet(keys [][]byte) ([][]byte, []bool, error) {
	var vals = make([][]byte, len(keys))
	var oks = make([]bool, len(keys))

	for i, k := range keys {
		v, err := s.db.Get(k)
		if err == nil {
			vals[i] = []byte(v)
			oks[i] = true
		} else {
			return vals, oks, err
		}
	}

	return vals, oks, nil
}

func (s *sniperStore) Set(key, value []byte) error {
	return s.db.Set(key, value, 0)
}

func (s *sniperStore) Get(key []byte) ([]byte, bool, error) {
	v, err := s.db.Get(key)
	return v, v != nil, err
}

func (s *sniperStore) Del(key []byte) (bool, error) {
	return s.db.Delete(key)
}

func (s *sniperStore) Keys(pattern []byte, limit int, withvalues bool) ([][]byte, [][]byte, error) {

	return nil, nil, nil
}

func (s *sniperStore) FlushDB() error {
	err := s.db.Close()
	// os.RemoveAll(s.path)
	return err
}
