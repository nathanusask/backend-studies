package datastore

import "context"

//go:generate mockgen -destination mock.go -source interface.go -package datastore
type Interface interface {
	Get(ctx context.Context, key interface{}) (interface{}, error)
	Upsert(ctx context.Context, key interface{}, value interface{}) error
	Delete(ctx context.Context, key interface{}) error
}
