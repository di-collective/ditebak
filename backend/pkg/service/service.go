package service

import "context"

// Service level abstraction
type Service interface {
	Reader
	Writer
}

// Reader abstraction to service layer
type Reader interface {
	// Get one
	Get(ctx context.Context, id string) (interface{}, error)

	// Find multiple
	Find(ctx context.Context, page, size int, opt map[string]interface{}) (total int64, rows []interface{}, err error)
}

// Writer abstraction to service layer
type Writer interface {
	// Create a new object
	Create(ctx context.Context, obj interface{}) (interface{}, error)

	// Update an existing object
	Update(ctx context.Context, id string, obj interface{}) (interface{}, error)

	// Delete an existing object virtually
	Delete(ctx context.Context, id string) error

	// Remove an existing object physically
	Remove(ctx context.Context, id string) error
}
