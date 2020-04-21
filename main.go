package main

import (
	"github.com/nasermirzaei89/github-visit-counter/app"
	"github.com/nasermirzaei89/github-visit-counter/repositories/boltdb"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
	"net/http"
	"os"
)

func main() {
	db, err := bolt.Open("db.bolt", 0666, nil)
	if err != nil {
		panic(errors.Wrap(err, "error on open db"))
	}
	defer func() { _ = db.Close() }()

	repo := boltdb.NewRepository(db)

	err = http.ListenAndServe(":"+os.Getenv("PORT"), app.NewHandler(repo))
	if err != nil {
		panic(errors.Wrap(err, "error on listen and serve"))
	}
}
