package main

import "encoding/json"

type CacheStorage struct {
	data map[string]interface{}
}

func NewCacheStorage() *CacheStorage {
	return &CacheStorage{
		data: make(map[string]interface{}),
	}
}

func (c *CacheStorage) Set(key, value []byte) error {
	c.data[string(key)] = value
	return nil
}

func (c *CacheStorage) Get(key []byte) ([]byte, bool, error) {
	value, wasFound := c.data[string(key)]
	var jsonValue []byte
	var err error

	if wasFound {
		jsonValue, err = json.Marshal(value)
	} else {
		jsonValue, err = nil, nil
	}

	return jsonValue, wasFound, err
}

func (c *CacheStorage) Delete(key []byte) error {
	delete(c.data, string(key))
	return nil
}

func (c *CacheStorage) Close() error {
	c.data = make(map[string]interface{})
	return nil
}
