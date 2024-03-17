package commands

import (
	"context"
	"payment-service/internal/modules/payment"
	"payment-service/internal/modules/payment/models/entity"
	"payment-service/internal/modules/payment/models/request"
	"payment-service/internal/pkg/databases/mongodb"
	"payment-service/internal/pkg/log"
	"time"

	wrapper "payment-service/internal/pkg/helpers"

	"go.mongodb.org/mongo-driver/bson"
)

type commandMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewCommandMongodbRepository(mongodb mongodb.Collections, log log.Logger) payment.MongodbRepositoryCommand {
	return &commandMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (c commandMongodbRepository) InsertOneOrder(ctx context.Context, order entity.Order) <-chan wrapper.Result {
	output := make(chan wrapper.Result)

	go func() {
		resp := <-c.mongoDb.InsertOne(mongodb.InsertOne{
			CollectionName: "order",
			Document:       order,
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (c commandMongodbRepository) InsertOnePayment(ctx context.Context, payment entity.PaymentHistory) <-chan wrapper.Result {
	output := make(chan wrapper.Result)

	go func() {
		resp := <-c.mongoDb.InsertOne(mongodb.InsertOne{
			CollectionName: "payment-history",
			Document:       payment,
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (c commandMongodbRepository) UpdatePaymentStatus(ctx context.Context, payload request.UpdatePaymentStatusReq) <-chan wrapper.Result {
	output := make(chan wrapper.Result)

	go func() {
		resp := <-c.mongoDb.UpdateOne(mongodb.UpdateOne{
			CollectionName: "payment-history",
			Filter: bson.M{
				"paymentId": payload.PaymentId,
			},
			Document: bson.M{
				"payment.transactionStatus": payload.TransactionStatus,
				"updatedAt":                 time.Now(),
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}
