package main

import (
	"time"
)

type Storage interface {
	Set(key string, value string, expiryCmd string, expiryValue int64)
	Get(key string) (string, bool)
}

type InMemoryStorage struct {
	data map[string]string;
	expiry map[string]int64;
}

func NewStorage() Storage {
	return &InMemoryStorage{
		data: make(map[string]string),
		expiry: make(map[string]int64),
	}
}

func (storage *InMemoryStorage) Set(key string, value string, expiryCmd string, expiryValue int64) {
	storage.data[key] = value

	switch expiryCmd {
	case "px":
		storage.expiry[key] = time.Now().UnixMilli() + expiryValue
	case "ex":
		storage.expiry[key] = time.Now().UnixMilli() + expiryValue * 1e3
	}
}

func (storage *InMemoryStorage) Get(key string) (value string, ok bool) {
	expiry, hasExpiry := storage.expiry[key]

	ok = false
	if (!hasExpiry) || (hasExpiry && expiry >= time.Now().UnixMilli()) {
		value = storage.data[key]
		ok = true
	}

	return value, ok
}