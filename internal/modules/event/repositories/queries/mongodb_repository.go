package queries

import (
	"context"
	"payment-service/internal/modules/event"
	"payment-service/internal/modules/event/models/entity"
	"payment-service/internal/pkg/databases/mongodb"
	wrapper "payment-service/internal/pkg/helpers"
	"payment-service/internal/pkg/log"

	"go.mongodb.org/mongo-driver/bson"
)

type queryMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewQueryMongodbRepository(mongodb mongodb.Collections, log log.Logger) event.MongodbRepositoryQuery {
	return &queryMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (q queryMongodbRepository) FindEventById(ctx context.Context, id string) <-chan wrapper.Result {
	var event entity.Event
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &event,
			CollectionName: "event",
			Filter: bson.M{
				"eventId": id,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}
