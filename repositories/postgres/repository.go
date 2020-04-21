package postgres

import (
	"context"
	"database/sql"
	"github.com/nasermirzaei89/github-visit-counter/app"
	"github.com/pkg/errors"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) app.Repository {
	repo := repository{db: db}

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS visits
(
    "key"    TEXT   NOT NULL PRIMARY KEY,
    "count" BIGINT NOT NULL DEFAULT 0
);
`)
	if err != nil {
		panic(errors.Wrap(err, "error on exec create table if not exists"))
	}

	return &repo
}

func (repo *repository) Get(ctx context.Context, key string) (int64, error) {
	var count int64

	err := repo.db.
		QueryRowContext(ctx, `SELECT "count" FROM visits WHERE "key" = $1;`, key).
		Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, errors.Wrap(err, "error on query select")
	}

	return count, nil
}

func (repo *repository) Visit(ctx context.Context, key string) (int64, error) {
	res, err := repo.Get(ctx, key)
	if err != nil {
		return 0, err
	}

	if res == 0 {
		_, err := repo.db.
			ExecContext(ctx, `INSERT INTO visits ("key", "count") SELECT $1, 1 WHERE NOT EXISTS (SELECT * FROM visits WHERE "key" = $1);`, key)
		if err != nil {
			return 0, errors.Wrap(err, "error on exec update")
		}
	} else {
		_, err := repo.db.
			ExecContext(ctx, `UPDATE visits SET "count" = "count" + 1 WHERE "key" = $1;`, key)
		if err != nil {
			return 0, errors.Wrap(err, "error on exec update")
		}
	}

	return res + 1, nil
}
