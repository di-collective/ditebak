package topic

import (
	"time"

	"github.com/di-collective/ditebak/backend/internal/domain/topic/dao"
	"github.com/di-collective/ditebak/backend/internal/rest/topic/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type delegate struct{}

func (del *delegate) Constructor() interface{} {
	return &dao.Topic{}
}

func (del *delegate) WillCreate(data interface{}) {
	now := time.Now()
	topic := data.(*dao.Topic)
	topic.CreatedAt = &now
	topic.State = dao.TopicStates.Draft()
}

func (del *delegate) DidCreate(created interface{}, id primitive.ObjectID) {
	topic := created.(*dao.Topic)
	topic.ID = id
}

func (del *delegate) WillUpdate(data interface{}, opt *options.UpdateOptions) {
	now := time.Now()
	topic := data.(*dto.Topic)
	topic.UpdatedAt = &now

	opt.SetUpsert(true)
}

func (del *delegate) DidUpdate(data interface{}, upsert *primitive.ObjectID) {
	if upsert != nil {
		topic := data.(*dto.Topic)
		topic.ID = *upsert
	}
}
