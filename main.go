package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/nasermirzaei89/github-visit-counter/app"
	"github.com/nasermirzaei89/github-visit-counter/repositories/postgres"
	"github.com/pkg/errors"
	"net/http"
	"os"
)

func main() {
	//db, err := bolt.Open("db.bolt", 0666, nil)
	//if err != nil {
	//	panic(errors.Wrap(err, "error on open db"))
	//}
	//defer func() { _ = db.Close() }()
	//
	//repo := boltdb.NewRepository(db)

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(errors.Wrap(err, "error on open postgres db"))
	}
	defer func() { _ = db.Close() }()

	err = db.Ping()
	if err != nil {
		panic(errors.Wrap(err, "error on ping postgres db"))
	}

	repo := postgres.NewRepository(db)

	err = http.ListenAndServe(":"+os.Getenv("PORT"), app.NewHandler(repo))
	if err != nil {
		panic(errors.Wrap(err, "error on listen and serve"))
	}
}
