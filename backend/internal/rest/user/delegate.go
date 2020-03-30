package user

import (
	"time"

	"github.com/di-collective/ditebak/backend/internal/domain/user/dao"
	"github.com/di-collective/ditebak/backend/internal/rest/user/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type delegate struct{}

func (del *delegate) Constructor() interface{} {
	return &dao.User{
		Reputation: 0,
	}
}

func (del *delegate) WillCreate(data interface{}) {
	now := time.Now()
	user := data.(*dao.User)
	user.CreatedAt = &now
}

func (del *delegate) DidCreate(created interface{}, id primitive.ObjectID) {
	user := created.(*dao.User)
	user.ID = id
}

func (del *delegate) WillUpdate(data interface{}, opt *options.UpdateOptions) {
	now := time.Now()
	user := data.(*dto.User)
	user.UpdatedAt = &now

	opt.SetUpsert(true)
}

func (del *delegate) DidUpdate(data interface{}, upsert *primitive.ObjectID) {
	if upsert != nil {
		user := data.(*dto.User)
		user.ID = *upsert
	}
}
