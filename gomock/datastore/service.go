package datastore

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

type row struct {
	value        interface{}
	created      time.Time
	lastModified time.Time
}

type service struct {
	db map[interface{}]*row
}

func New(db map[interface{}]*row) Interface {
	return &service{
		db: db,
	}
}

func (s *service) Get(ctx context.Context, key interface{}) (interface{}, error) {
	if row, ok := s.db[key]; ok {
		return row.value, nil
	}
	return nil, fmt.Errorf("cannot find the associated value with key: %+v", key)
}

func (s *service) Upsert(ctx context.Context, key interface{}, value interface{}) error {
	if v, err := s.Get(ctx, key); err == nil {
		if !reflect.DeepEqual(v, value) {
			// the value is newer
			r := s.db[key]
			r.value = value
			r.lastModified = time.Now()
			return nil
		}
		return fmt.Errorf("record already exists")
	}
	r := &row{
		value:        value,
		created:      time.Now(),
		lastModified: time.Now(),
	}
	s.db[key] = r
	return nil
}

func (s *service) Delete(ctx context.Context, key interface{}) error {
	if _, err := s.Get(ctx, key); err != nil {
		return fmt.Errorf("no record is found in regards to the key %+v", key)
	}
	delete(s.db, key)
	return nil
}
