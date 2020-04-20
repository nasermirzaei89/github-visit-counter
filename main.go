package main

import (
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
	"net/http"
	"os"
)

func main() {
	db, err := bolt.Open(env("DB_PATH", "db.bolt"), 0666, nil)
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

	err = http.ListenAndServe(env("API_ADDRESS", ":80"), NewHandler(db))
	if err != nil {
		panic(errors.Wrap(err, "error on listen and serve"))
	}
}

func env(key, def string) string {
	res, ok := os.LookupEnv(key)
	if ok {
		return res
	}

	return def
}
