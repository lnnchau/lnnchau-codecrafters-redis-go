package main

type Storage interface {
	Set(key string, value string)
	Get(key string) string
}

type InMemoryStorage struct {
	data map[string]string;
	expiry map[string]float64;
}

func NewStorage() Storage {
	return &InMemoryStorage{
		data: make(map[string]string),
		expiry: make(map[string]float64),
	}
}

func (storage *InMemoryStorage) Set(key string, value string) {
	storage.data[key] = value
}

func (storage *InMemoryStorage) Get(key string) string {
	return storage.data[key]
}