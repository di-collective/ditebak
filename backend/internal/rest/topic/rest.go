package topic

import (
	"context"
	"reflect"

	log "github.com/sirupsen/logrus"

	"github.com/di-collective/ditebak/backend/internal/rest/topic/dto"
	"github.com/di-collective/ditebak/backend/pkg/queryables"
	"github.com/di-collective/ditebak/backend/pkg/repo/mongorepo"
	"github.com/di-collective/ditebak/backend/pkg/rest"
	"github.com/di-collective/ditebak/backend/pkg/service/basic"
	"go.mongodb.org/mongo-driver/mongo"
)

// New instance of Topic REST API
func New(coll *mongo.Collection) rest.REST {
	delegate := &delegate{}
	return rest.New(&rest.Config{
		Resource: "topics",
		Service: basic.New(mongorepo.New(
			/* collection   */ coll,
			/* default sort */ map[string]int{"created_at": -1},
			/* constructor  */ delegate.Constructor,
			/* id assigner  */ delegate)),
		CreatePayload: delegate.Constructor, //dto = dao
		UpdatePayload: func() interface{} {
			//uses dto.Topic to allow partial update
			return &dto.Topic{}
		},
		Queryables: queryables.Collection{
			{DtoKey: "state", DaoKey: "state", TypeOf: reflect.Array,
				// (state=a,b,c) --> state: {$in: [a, b, c]}
				Transform: func(key string, value interface{}) (string, interface{}) {
					return key, map[string]interface{}{
						"$in": value,
					}
				}},
		},
	})
}

// Teardown REST API
func Teardown(ctx context.Context) error {
	// TODO: graceful shutdown
	log.Info("SHUTTING DOWN")
	return nil
}
