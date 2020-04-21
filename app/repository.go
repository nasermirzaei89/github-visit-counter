package app

import "context"

type Repository interface {
	Get(ctx context.Context, key string) (count int64, err error)
	Visit(ctx context.Context, key string) (count int64, err error)
}
