package boltdb

import (
	"context"
	"github.com/nasermirzaei89/github-visit-counter/app"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
	"strconv"
)

type repository struct {
	db *bolt.DB
}

func NewRepository(db *bolt.DB) app.Repository {
	repo := repository{db: db}

	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("visits"))
		if err != nil {
			return errors.Wrap(err, "error on create bucket")
		}
		return nil
	})
	if err != nil {
		panic(errors.Wrap(err, "error on update transaction"))
	}

	return &repo
}

func (repo *repository) Get(ctx context.Context, key string) (int64, error) {
	var count int64

	err := repo.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("visits"))
		v := b.Get([]byte(key))
		if v == nil {
			v = []byte("0")
		}

		count, _ = strconv.ParseInt(string(v), 10, 64)

		return nil
	})

	if err != nil {
		return 0, errors.Wrap(err, "error on update transaction")
	}

	return count, nil
}

func (repo *repository) Visit(ctx context.Context, key string) (int64, error) {
	var count int64

	err := repo.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("visits"))
		v := b.Get([]byte(key))
		if v == nil {
			v = []byte("0")
		}

		count, _ = strconv.ParseInt(string(v), 10, 64)

		count++

		err := b.Put([]byte(key), []byte(strconv.FormatInt(count, 10)))
		if err != nil {
			return errors.Wrap(err, "error on put in transaction")
		}

		return nil
	})

	if err != nil {
		return 0, errors.Wrap(err, "error on update transaction")
	}

	return count, nil
}
