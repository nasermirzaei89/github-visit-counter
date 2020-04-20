package main

import (
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

func main() {
	db, err := bolt.Open("db.bolt", 0666, nil)
	if err != nil {
		panic(errors.Wrap(err, "error on open db"))
	}
	defer func() { _ = db.Close() }()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("visits"))
		if err != nil {
			return errors.Wrap(err, "error on create bucket")
		}
		return nil
	})
	if err != nil {
		panic(errors.Wrap(err, "error on update transaction"))
	}

	err = http.ListenAndServe(":9898", NewHandler(db))
	if err != nil {
		panic(errors.Wrap(err, "error on listen and serve"))
	}
}
