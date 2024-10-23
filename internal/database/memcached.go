package database

import (
	"log"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemCacheDB struct {
	DB *memcache.Client
}

func NewMemCacheDB() MemCacheDB {
	mc := memcache.New(":11211")
	if mc == nil {
		log.Fatal("MEMCACHE CLIENT IS NIL")
	}

	return MemCacheDB{DB: mc}
}
