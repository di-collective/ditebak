package bet

import (
	"time"

	"github.com/di-collective/ditebak/backend/internal/domain/bet/dao"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type delegate struct{}

func (del *delegate) Constructor() interface{} {
	return &dao.Bet{}
}

func (del *delegate) WillCreate(data interface{}) {
	now := time.Now()
	bet := data.(*dao.Bet)
	bet.CreatedAt = &now
	bet.State = dao.BetStates.Placed()
}

func (del *delegate) DidCreate(created interface{}, id primitive.ObjectID) {
	bet := created.(*dao.Bet)
	bet.ID = id
}

func (del *delegate) WillUpdate(data interface{}, opt *options.UpdateOptions) {
	now := time.Now()
	bet := data.(*dao.Bet)
	bet.UpdatedAt = &now

	opt.SetUpsert(true)
}

func (del *delegate) DidUpdate(data interface{}, upsert *primitive.ObjectID) {
	if upsert != nil {
		bet := data.(*dao.Bet)
		bet.ID = *upsert
	}
}
