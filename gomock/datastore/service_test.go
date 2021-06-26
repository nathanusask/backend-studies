package datastore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_service_Delete(t *testing.T) {
	tests := []struct {
		name        string
		db          map[interface{}]*row
		key         interface{}
		expectedErr error
	}{
		{
			name: "Happy path",
			db: map[interface{}]*row{
				"foo": {
					value:        "bar",
					created:      time.Date(2021, 5, 1, 12, 10, 10, 0, time.UTC),
					lastModified: time.Date(2021, 5, 1, 12, 10, 10, 0, time.UTC),
				},
			},
			key:         "foo",
			expectedErr: nil,
		},
		{
			name:        "fails when the key is missing from the db",
			db:          nil,
			key:         "foo",
			expectedErr: fmt.Errorf("no record is found in regards to the key %+v", "foo"),
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.db)

			err := s.Delete(ctx, tt.key)
			if tt.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err, tt.expectedErr.Error())
			}
		})
	}
}

func Test_service_Get(t *testing.T) {
	tests := []struct {
		name          string
		db            map[interface{}]*row
		key           interface{}
		expectedValue interface{}
		wantErr       bool
		expectedErr   error
	}{
		{
			name: "Happy path",
			db: map[interface{}]*row{
				"foo": {
					value:        "bar",
					created:      time.Date(2021, 5, 1, 12, 10, 10, 0, time.UTC),
					lastModified: time.Date(2021, 5, 1, 12, 10, 10, 0, time.UTC),
				},
			},
			key:           "foo",
			expectedValue: "bar",
			wantErr:       false,
			expectedErr:   nil,
		},
		{
			name:          "fails when the key is not found in the db",
			db:            nil,
			key:           "foo",
			expectedValue: nil,
			wantErr:       true,
			expectedErr:   fmt.Errorf("cannot find the associated value with key: %+v", "foo"),
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.db)
			got, err := s.Get(ctx, tt.key)
			if !tt.wantErr {
				assert.Equal(t, got, tt.expectedValue)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err, tt.expectedErr.Error())
			}
		})
	}
}

func Test_service_Upsert(t *testing.T) {
	tests := []struct {
		name          string
		db            map[interface{}]*row
		key           interface{}
		value         interface{}
		wantErr       bool
		expectedError error
	}{
		{
			name:          "success when the db is empty",
			db:            map[interface{}]*row{},
			key:           "foo",
			value:         "bar",
			wantErr:       false,
			expectedError: nil,
		},
		{
			name: "success when the record is newer",
			db: map[interface{}]*row{
				"foo": {
					value:        "ba",
					created:      time.Date(2021, 5, 1, 12, 10, 10, 0, time.UTC),
					lastModified: time.Date(2021, 5, 1, 12, 10, 10, 0, time.UTC),
				},
			},
			key:           "foo",
			value:         "bar",
			wantErr:       false,
			expectedError: nil,
		}, {
			name: "fails when the record already exists",
			db: map[interface{}]*row{
				"foo": {
					value:        "bar",
					created:      time.Date(2021, 5, 1, 12, 10, 10, 0, time.UTC),
					lastModified: time.Date(2021, 5, 1, 12, 10, 10, 0, time.UTC),
				},
			},
			key:           "foo",
			value:         "bar",
			wantErr:       true,
			expectedError: fmt.Errorf("record already exists"),
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.db)
			err := s.Upsert(ctx, tt.key, tt.value)
			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Errorf(t, err, tt.expectedError.Error())
			}
		})
	}
}
