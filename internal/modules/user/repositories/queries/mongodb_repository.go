package queries

import (
	"context"
	user "payment-service/internal/modules/user"
	userEntity "payment-service/internal/modules/user/models/entity"
	"payment-service/internal/pkg/databases/mongodb"
	wrapper "payment-service/internal/pkg/helpers"
	"payment-service/internal/pkg/log"

	"go.mongodb.org/mongo-driver/bson"
)

type queryMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewQueryMongodbRepository(mongodb mongodb.Collections, log log.Logger) user.MongodbRepositoryQuery {
	return &queryMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (q queryMongodbRepository) FindOneUserId(ctx context.Context, userId string) <-chan wrapper.Result {
	var user userEntity.User
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &user,
			CollectionName: "users",
			Filter: bson.M{
				"userId": userId,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}
