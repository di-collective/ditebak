package basic

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/di-collective/ditebak/backend/pkg/exception"
	"github.com/di-collective/ditebak/backend/pkg/repo"
)

// Service abstraction
type Service struct {
	rps repo.Repository
}

// New basic service
// @rps: persistence repository
func New(rps repo.Repository) *Service {
	return &Service{
		rps: rps,
	}
}

// Get one
func (svc *Service) Get(ctx context.Context, id string) (interface{}, error) {
	res, err := svc.rps.Get(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return res, exception.New(http.StatusNotFound, "Resource with ID: %s, is not found", id)
		}
	}

	return res, err
}

// Find multiple
func (svc *Service) Find(ctx context.Context, page, size int, opt map[string]interface{}) (total int64, rows []interface{}, err error) {
	fo := repo.FindOptions{
		Page:           page,
		Size:           size,
		IncludeRemoved: false,
		Params:         opt,
	}

	return svc.rps.Find(ctx, fo)
}

// Create a new object
func (svc *Service) Create(ctx context.Context, obj interface{}) (interface{}, error) {
	err := svc.rps.Create(ctx, obj)
	if mwe, ok := err.(mongo.WriteException); ok {
		if len(mwe.WriteErrors) > 0 {
			e := mwe.WriteErrors[0]
			if e.Code == 11000 {
				return nil, exception.New(http.StatusConflict, "Duplicate resource already exists")
			}
		}
		return nil, err
	}

	return obj, err
}

// Update an existing object
func (svc *Service) Update(ctx context.Context, id string, obj interface{}) (interface{}, error) {
	err := svc.rps.Update(ctx, id, obj)
	return obj, err
}

// Delete an existing object virtually
func (svc *Service) Delete(ctx context.Context, id string) error {
	return svc.rps.Delete(ctx, id)
}

// Remove an existing object physically
func (svc *Service) Remove(ctx context.Context, id string) error {
	return svc.rps.Remove(ctx, id)
}
