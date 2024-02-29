package queries

import (
	"context"
	"payment-service/internal/modules/ticket"
	"payment-service/internal/modules/ticket/models/entity"
	"payment-service/internal/pkg/databases/mongodb"
	wrapper "payment-service/internal/pkg/helpers"
	"payment-service/internal/pkg/log"

	"go.mongodb.org/mongo-driver/bson"
)

type queryMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewQueryMongodbRepository(mongodb mongodb.Collections, log log.Logger) ticket.MongodbRepositoryQuery {
	return &queryMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (q queryMongodbRepository) FindBankTicketByTicketNumber(ctx context.Context, ticketNumber string, eventId string) <-chan wrapper.Result {
	var ticket entity.BankTicket
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &ticket,
			CollectionName: "bank-ticket",
			Filter: bson.M{
				"ticketNumber": ticketNumber,
				"eventId":      eventId,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}
