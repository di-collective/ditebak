package mongorepo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Event delegates
type Event interface {
	WillCreate(interface{})
	DidCreate(interface{}, primitive.ObjectID)

	WillUpdate(interface{}, *options.UpdateOptions)

	// @obj: is the object being updated
	// @upsert: is the ID of object if upsert is done. nil if no upsert
	DidUpdate(obj interface{}, upsert *primitive.ObjectID)
}
