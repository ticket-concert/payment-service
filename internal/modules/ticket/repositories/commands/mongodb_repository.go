package commands

import (
	"context"
	"payment-service/internal/modules/ticket"
	"payment-service/internal/modules/ticket/models/request"
	"payment-service/internal/pkg/databases/mongodb"
	wrapper "payment-service/internal/pkg/helpers"
	"payment-service/internal/pkg/log"

	"go.mongodb.org/mongo-driver/bson"
)

type commandMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewCommandMongodbRepository(mongodb mongodb.Collections, log log.Logger) ticket.MongodbRepositoryCommand {
	return &commandMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (c commandMongodbRepository) UpdateBankTicket(ctx context.Context, payload request.UpdateBankTicketReq) <-chan wrapper.Result {
	output := make(chan wrapper.Result)

	go func() {
		resp := <-c.mongoDb.UpdateOne(mongodb.UpdateOne{
			CollectionName: "bank-ticket",
			Filter: bson.M{
				"ticketNumber": payload.TicketNumber,
			},
			Document: bson.M{
				"paymentStatus": payload.PaymentStatus,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}
