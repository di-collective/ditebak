package credential

import (
	"github.com/di-collective/ditebak/backend/internal/domain/credential/dao"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type delegate struct{}

func (del *delegate) Constructor() interface{} {
	return &dao.Credential{}
}

func (del *delegate) WillCreate(data interface{}) {
	// do nothing
}

func (del *delegate) DidCreate(created interface{}, id primitive.ObjectID) {
	cred := created.(*dao.Credential)
	cred.ID = id
}

func (del *delegate) WillUpdate(data interface{}, opt *options.UpdateOptions) {
	opt.SetUpsert(true)
}

func (del *delegate) DidUpdate(data interface{}, upsert *primitive.ObjectID) {
	if upsert != nil {
		cred := data.(*dao.Credential)
		cred.ID = *upsert
	}
}
