package main

import "github.com/cockroachdb/pebble"

type Storage interface {
	Get(key []byte) ([]byte, bool, error)
	Set(key []byte, val []byte) error
	Delete(key []byte) error
	Close() error
}

type PebbleStorage struct {
	db *pebble.DB
}

func NewPebbleStorage(dirname string) *PebbleStorage {
	db, err := pebble.Open(dirname, &pebble.Options{})
	if err != nil {
		panic(err)
	}
	return &PebbleStorage{
		db: db,
	}
}

func (p *PebbleStorage) Set(key, val []byte) error {
	return p.db.Set(key, val, pebble.Sync)
}

func (p *PebbleStorage) Get(key []byte) ([]byte, bool, error) {
	value, closer, err := p.db.Get(key)
	wasFound := true

	if err == pebble.ErrNotFound {
		wasFound = false
		err = nil
	}

	if closer != nil {
		closer.Close()
	}

	return value, wasFound, err
}

func (p *PebbleStorage) Delete(key []byte) error {
	return p.db.Delete(key, pebble.Sync)
}

func (p *PebbleStorage) Close() error {
	return p.db.Close()
}
