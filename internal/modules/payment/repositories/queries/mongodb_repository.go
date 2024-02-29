package queries

import (
	"context"
	"payment-service/internal/modules/payment"
	"payment-service/internal/modules/payment/models/entity"
	"payment-service/internal/modules/payment/models/request"
	"payment-service/internal/pkg/databases/mongodb"
	wrapper "payment-service/internal/pkg/helpers"
	"payment-service/internal/pkg/log"

	"go.mongodb.org/mongo-driver/bson"
)

type queryMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewQueryMongodbRepository(mongodb mongodb.Collections, log log.Logger) payment.MongodbRepositoryQuery {
	return &queryMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (q queryMongodbRepository) FindOrderByTicket(ctx context.Context, ticketNumber string) <-chan wrapper.Result {
	var order entity.Order
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &order,
			CollectionName: "order",
			Filter: bson.M{
				"ticketNumber": ticketNumber,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindPaymentById(ctx context.Context, id string) <-chan wrapper.Result {
	var payment entity.PaymentHistory
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &payment,
			CollectionName: "payment-history",
			Filter: bson.M{
				"paymentId":      id,
				"isValidPayment": true,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindPaymentByTicketNumber(ctx context.Context, ticketNumber string) <-chan wrapper.Result {
	var payment entity.PaymentHistory
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &payment,
			CollectionName: "payment-history",
			Filter: bson.M{
				"ticketNumber":   ticketNumber,
				"isValidPayment": true,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindPaymentByUser(ctx context.Context, payload request.PaymentList) <-chan wrapper.Result {
	var payments []entity.PaymentHistory
	var countData int64
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindAllData(mongodb.FindAllData{
			Result:         &payments,
			CountData:      &countData,
			CollectionName: "payment-history",
			Filter:         bson.M{"userId": payload.UserId},
			Sort: &mongodb.Sort{
				FieldName: "createdAt",
				By:        mongodb.SortDescending,
			},
			Page: payload.Page,
			Size: payload.Size,
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindPaymentStatusById(ctx context.Context, id string) <-chan wrapper.Result {
	var payment entity.PaymentHistory
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &payment,
			CollectionName: "payment-history",
			Filter: bson.M{
				"paymentId": id,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}
