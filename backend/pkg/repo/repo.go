package repo

import (
	"context"
)

// Repository level abstraction
type Repository interface {
	Reader
	Writer
}

// Reader abstraction to persistent layer
type Reader interface {
	// Get one
	Get(ctx context.Context, id string) (interface{}, error)

	// Find multiple
	Find(ctx context.Context, opt FindOptions) (total int64, rows []interface{}, err error)
}

// Writer abstraction to persistent layer
type Writer interface {
	// Create a new object
	Create(ctx context.Context, obj interface{}) error

	// Update an existing object
	Update(ctx context.Context, id string, obj interface{}) error

	// Delete an existing object virtually
	Delete(ctx context.Context, id string) error

	// Remove an existing object physically
	Remove(ctx context.Context, id string) error
}
