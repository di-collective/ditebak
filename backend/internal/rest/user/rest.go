package user

import (
	"context"
	"reflect"

	log "github.com/sirupsen/logrus"

	"github.com/di-collective/ditebak/backend/internal/rest/user/dto"
	"github.com/di-collective/ditebak/backend/pkg/queryables"
	"github.com/di-collective/ditebak/backend/pkg/repo/mongorepo"
	"github.com/di-collective/ditebak/backend/pkg/rest"
	"github.com/di-collective/ditebak/backend/pkg/service/basic"
	"go.mongodb.org/mongo-driver/mongo"
)

// New instance of User REST API
func New(coll *mongo.Collection) rest.REST {
	delegate := &delegate{}
	return rest.New(&rest.Config{
		Resource: "users",
		Service: basic.New(mongorepo.New(
			/* collection   */ coll,
			/* default sort */ map[string]int{"created_at": -1},
			/* constructor  */ delegate.Constructor,
			/* id assigner  */ delegate)),
		CreatePayload: delegate.Constructor,
		UpdatePayload: func() interface{} {
			return &dto.User{}
		},
		Convert: nil, // dto == dao
		Queryables: queryables.Collection{
			{DtoKey: "provider", DaoKey: "provider", TypeOf: reflect.String},
			{DtoKey: "email", DaoKey: "email", TypeOf: reflect.String},
		},
	})
}

// Teardown REST API
func Teardown(ctx context.Context) error {
	// TODO: graceful shutdown
	log.Info("SHUTTING DOWN")
	return nil
}
