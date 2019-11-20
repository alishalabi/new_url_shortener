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

// Set function to open a new db connection
func openDatabase(stumb string) *bbolt.DB {
  // Open db file in current directory OR create if does not exist
  db, err := bbolt.Open(stumb, 0600, nil)
  if err!= nil {
    Panic(err)
  }

  // Create buckets
  var tables = [...][]byte {
    tableURLs,
  }

  db.Update(func(tx *bbolt.Tx) (err error) {
    for _, table := range tables {
      _, err = tx.CreateBucketIfNotExists(table)
      if err != nil {
        Panic(err)
      }
    })
    return
  }
  return db
}


// Set function to return new DB instance (connection is open)
// DB implements the Store
func NewDB(stumb string) *DB {
  return &DB{
    db: openDatabase(stumb),
  }
}

// Set shorten url and shorten url key
func (d *DB) Set(key string, value string) error {
  return d.db.Update(func(tx *bbolt.Tx) error {
    b, err := tx.CreateBucketIfNotExists(tableURLs)
    if err != nil {
      return err
    }
    k := []byte(key)
    valueB := []byte(value)
    c := b.Cursor()

    found := false
    valueB := []byte(value)
    for k, v := c.First();k != nil; k, v = c.Next() {
      if bytes.Equal(valueB, v) {
        found = true
        break
      }
    }
    // If value exists, end function
    if found {
      return nil
    }

    return b.Put(k, []byte(value))
  })
}

// Create "clear" method to easily clear all database entries
func (d *DB) Clear() error {
  return d.db.Update(func(tx *bbolt.Tx) error {
    return tx.DeleteBucket(tableURLs)
    })
}

// Create "get" method to return url by its key
// If not found, return empty string
func (d *DB) Get(key string) (value string) {
  keyB := []byte(key)
  d.db.Update(func(tx *bbolt.Tx) error {
    b := tx.Bucket(tableURLs)
    if b == nil {
      return nil
    }
    c:= b.Cursor()
    for k, v := c.First(); k != nil; k, v = c.Next() {
      if bytes.Equal(keyB, k) {
        value = string(v)
        break
      }
    }
    return nil
  })

  return
}

// Create "GetByValue" method returns all keys
// for a specific (original) url value.
func (d *DB) GetByValue(value string) (keys []strings) {
  valueB := []byte(value)
  d.db.Update(func (tx *bbolt.Tx) error {
    b := tx.Bucket(tableURLs)
    if b == nil {
      return nil
    }
    c := b.Cursor()
    for k, v := c.First(); k != nil; k, v = c.Next() {
      if bytes.Equal(valueB, v) {
        keys = append(keys, string(k))
      }
    }
    return nil
  })
  return 
}
