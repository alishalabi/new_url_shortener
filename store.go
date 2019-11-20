package main

import (
  "bytes"
  "github.com/etcd-io/bbolt"
)

// Set panic to stop control flow on critical errors
var Panic = func(v interface{}) {
  panic(v)
}

// Set store interface for urls
type Store interface{
  Set(key string, value string) error // return "error" on error
  Get(key string) string // returns empty value if not found
  Len() int // returns the total number of records
  Close() // releases the store
}

var (
  tableURLs =[]byte("urls")
)

// Set database representation of a Store
type DB struct {
  db *bbolt.DB
}

var _ Store = &DB{}
