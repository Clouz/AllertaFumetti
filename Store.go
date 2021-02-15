package main

import (
	"log"

	"github.com/timshannon/bolthold"
	bolt "go.etcd.io/bbolt"
)

//https://github.com/timshannon/bolthold/blob/master/example_test.go
func data(i []item) {
	store, err := bolthold.Open("database", 0666, nil)
	defer store.Close()

	if err != nil {
		// handle error
		log.Fatal(err)
	}

	err = store.Bolt().Update(func(tx *bolt.Tx) error {
		for i := range data {
			err := store.TxInsert(tx, data[i].ID, data[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}
