package payment

import (
	"context"
	"payment-service/internal/modules/payment/models/entity"
	"payment-service/internal/modules/payment/models/request"
	"payment-service/internal/modules/payment/models/response"
	wrapper "payment-service/internal/pkg/helpers"
)

type UsecaseCommand interface {
	CreatePayment(origCtx context.Context, payload request.PaymentReq) (*response.PaymentResp, error)
	CreateTicketOrder(origCtx context.Context, payload request.TicketOrderReq) (*string, error)
}

type UsecaseQuery interface {
	FindPaymentStatus(origCtx context.Context, payload request.PaymentStatusReq) (*response.PaymentStatusResp, error)
	FindOrderPayment(origCtx context.Context, payload request.GetOrderPaymentReq) (*response.GetOrderPaymentResp, error)
	FindPaymentList(origCtx context.Context, payload request.PaymentList) (*response.OrderListResp, error)
}

type MongodbRepositoryQuery interface {
	FindOrderByTicket(ctx context.Context, ticketNumber string) <-chan wrapper.Result
	FindPaymentById(ctx context.Context, id string) <-chan wrapper.Result
	FindPaymentByTicketNumber(ctx context.Context, ticketNumber string) <-chan wrapper.Result
	FindPaymentByUser(ctx context.Context, payload request.PaymentList) <-chan wrapper.Result
	FindPaymentStatusById(ctx context.Context, id string) <-chan wrapper.Result
}

type MongodbRepositoryCommand interface {
	InsertOneOrder(ctx context.Context, order entity.Order) <-chan wrapper.Result
	InsertOnePayment(ctx context.Context, payment entity.PaymentHistory) <-chan wrapper.Result
	UpdatePaymentStatus(ctx context.Context, payload request.UpdatePaymentStatusReq) <-chan wrapper.Result
}

type MidtransRepositoryCommand interface {
	TransferBank(ctx context.Context, payload request.BankTransferRequest) (*response.BankTransferResponse, error)
}

type MidtransRepositoryQuery interface {
	GetTransactionStatus(ctx context.Context, transactionId string) (*response.TransactionStatusResponse, error)
}
