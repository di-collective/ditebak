package mongorepo

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/di-collective/ditebak/backend/pkg/repo"
)

var (
	virtualDelete = bson.M{"$set": bson.M{"_deleted": true}}
)

// Repo abstraction
type Repo struct {
	collection  *mongo.Collection
	sort        map[string]int
	constructor func() interface{}
	delegates   Event
}

// New Repo using mongodb
// @coll: mongo collection
// @sort: Please provide a default sort
// @con: Please provide a function to a constructor/factory which return a pointer to a struct
// @aid: Please proved a function to which ID must be assigned
func New(coll *mongo.Collection,
	sort map[string]int,
	con func() interface{},
	del Event) *Repo {
	return &Repo{
		collection:  coll,
		sort:        sort,
		constructor: con,
		delegates:   del,
	}
}

// Get one
func (r *Repo) Get(ctx context.Context, id string) (interface{}, error) {
	log.Traceln(r.collection.Name(), "GET", id)

	_id, _ := primitive.ObjectIDFromHex(id)
	res := r.collection.FindOne(ctx, bson.M{"_id": _id})
	dbo := r.constructor()
	err := res.Decode(dbo)

	return dbo, err
}

// Find multiple
func (r *Repo) Find(ctx context.Context, opt repo.FindOptions) (int64, []interface{}, error) {
	trace := fmt.Sprintf("%s %s", r.collection.Name(), "FIND")

	// 1. set paging
	fo := options.Find().
		SetSkip(int64(opt.Skip())).
		SetLimit(int64(opt.Size))

	// 2. set sort
	if opt.Sort != nil {
		fo.SetSort(bson.M{
			opt.Sort.Field: sortDirection(opt.Sort),
		})
	} else if r.sort != nil && len(r.sort) > 0 {
		fo.SetSort(r.sort)
	}

	// 3. set query
	fi := opt.Params
	if !opt.IncludeRemoved {
		fi["_deleted"] = map[string]bool{"$exists": false}
	}

	log.Traceln(trace, fi)
	cur, err := r.collection.Find(ctx, opt.Params, fo)
	if err != nil {
		return 0, nil, err
	}

	result := []interface{}{} // I don't want null slice, I want empty slice
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		dbo := r.constructor()
		if err = cur.Decode(dbo); err != nil {
			return 0, nil, err
		}

		result = append(result, dbo)
	}

	count, _ := r.collection.CountDocuments(ctx, opt.Params)

	log.Traceln(trace, "total:", count, "result:", result)
	return count, result, nil
}

// Create a new object
func (r *Repo) Create(ctx context.Context, obj interface{}) error {
	r.delegates.WillCreate(obj)

	res, err := r.collection.InsertOne(ctx, obj)
	if err != nil {
		return err
	}

	r.delegates.DidCreate(obj, res.InsertedID.(primitive.ObjectID))
	return nil
}

// Update an existing object
func (r *Repo) Update(ctx context.Context, id string, obj interface{}) error {
	uo := options.Update()
	r.delegates.WillUpdate(obj, uo)

	_id, _ := primitive.ObjectIDFromHex(id)
	setter := bson.M{"$set": obj}
	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": _id}, setter, uo)
	if err != nil {
		return err
	}

	var uid *primitive.ObjectID
	if res.UpsertedID != nil {
		pid, _ := res.UpsertedID.(primitive.ObjectID)
		uid = &pid
	}

	r.delegates.DidUpdate(obj, uid)
	return nil
}

// Delete an existing object virtually
func (r *Repo) Delete(ctx context.Context, id string) error {
	_id, _ := primitive.ObjectIDFromHex(id)
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": _id}, virtualDelete)
	if err != nil {
		return err
	}

	return nil
}

// Remove an existing object physically
func (r *Repo) Remove(ctx context.Context, id string) error {
	_id, _ := primitive.ObjectIDFromHex(id)
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": _id})
	if err != nil {
		return err
	}

	return nil
}

func sortDirection(opt *repo.SortOption) int {
	if opt.Descending {
		return -1
	}

	return 1
}
