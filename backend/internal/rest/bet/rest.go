package bet

import (
	"context"
	"reflect"

	log "github.com/sirupsen/logrus"

	"github.com/di-collective/ditebak/backend/pkg/queryables"
	"github.com/di-collective/ditebak/backend/pkg/repo/mongorepo"
	"github.com/di-collective/ditebak/backend/pkg/rest"
	"github.com/di-collective/ditebak/backend/pkg/service/basic"
	"go.mongodb.org/mongo-driver/mongo"
)

// New instance of Bet REST API
func New(coll *mongo.Collection) rest.REST {
	delegate := &delegate{}
	return rest.New(&rest.Config{
		Resource: "bets",
		Service: basic.New(mongorepo.New(
			/* collection    */ coll,
			/* default sort  */ map[string]int{"created_at": -1},
			/* constructor   */ delegate.Constructor,
			/* event handler */ delegate)),
		CreatePayload: delegate.Constructor,
		UpdatePayload: delegate.Constructor,
		Convert:       nil, // dto == dao
		Queryables: queryables.Collection{
			{DtoKey: "topic", DaoKey: "topic_id", TypeOf: reflect.String},
			{DtoKey: "owner", DaoKey: "owner", TypeOf: reflect.String},
			{DtoKey: "state", DaoKey: "state", TypeOf: reflect.String},
		},
	})
}

// Teardown REST API
func Teardown(ctx context.Context) error {
	// TODO: graceful shutdown
	log.Info("SHUTTING DOWN")
	return nil
}
