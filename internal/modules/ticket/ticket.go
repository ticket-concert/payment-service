package ticket

import (
	"context"
	"payment-service/internal/modules/ticket/models/request"
	"payment-service/internal/modules/ticket/models/response"
	wrapper "payment-service/internal/pkg/helpers"
)

type UsecaseQuery interface {
	FindTickets(origCtx context.Context, payload request.TicketReq) (*response.TicketResp, error)
	FindOnlineTicket(origCtx context.Context, payload request.TicketReq) (*response.Ticket, error)
	FindAvailableTicket(origCtx context.Context) ([]response.TicketCountry, error)
}

type MongodbRepositoryQuery interface {
	FindBankTicketByTicketNumber(ctx context.Context, ticketNumber string, eventId string) <-chan wrapper.Result
}

type MongodbRepositoryCommand interface {
	UpdateBankTicket(ctx context.Context, payload request.UpdateBankTicketReq) <-chan wrapper.Result
}
