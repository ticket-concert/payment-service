package ticket

import (
	"context"
	"payment-service/internal/modules/ticket/models/request"
	wrapper "payment-service/internal/pkg/helpers"
)

type MongodbRepositoryQuery interface {
	FindBankTicketByTicketNumber(ctx context.Context, ticketNumber string, eventId string) <-chan wrapper.Result
}

type MongodbRepositoryCommand interface {
	UpdateBankTicket(ctx context.Context, payload request.UpdateBankTicketReq) <-chan wrapper.Result
}
