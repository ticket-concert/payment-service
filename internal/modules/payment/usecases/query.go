package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"payment-service/internal/modules/event"
	"payment-service/internal/modules/payment"
	"payment-service/internal/modules/payment/models/entity"
	"payment-service/internal/modules/payment/models/request"
	"payment-service/internal/modules/payment/models/response"
	"payment-service/internal/modules/ticket"
	"payment-service/internal/pkg/constants"
	"payment-service/internal/pkg/errors"
	"payment-service/internal/pkg/helpers"
	"payment-service/internal/pkg/log"
	"payment-service/internal/pkg/redis"
	"time"

	eventEntity "payment-service/internal/modules/event/models/entity"
	ticketEntity "payment-service/internal/modules/ticket/models/entity"
	userDto "payment-service/internal/modules/user/models/dto"

	"go.elastic.co/apm"
)

type queryUsecase struct {
	paymentRepositoryQuery  payment.MongodbRepositoryQuery
	midtransRepositoryQuery payment.MidtransRepositoryQuery
	ticketRepositoryQuery   ticket.MongodbRepositoryQuery
	eventRepositoryQuery    event.MongodbRepositoryQuery
	logger                  log.Logger
	redis                   redis.Collections
}

func NewQueryUsecase(pmq payment.MongodbRepositoryQuery, mrq payment.MidtransRepositoryQuery, tmq ticket.MongodbRepositoryQuery, emq event.MongodbRepositoryQuery, log log.Logger, rc redis.Collections) payment.UsecaseQuery {
	return queryUsecase{
		paymentRepositoryQuery:  pmq,
		midtransRepositoryQuery: mrq,
		ticketRepositoryQuery:   tmq,
		eventRepositoryQuery:    emq,
		logger:                  log,
		redis:                   rc,
	}
}

func (q queryUsecase) FindPaymentStatus(origCtx context.Context, payload request.PaymentStatusReq) (*response.PaymentStatusResp, error) {
	domain := "paymentUsecase-FindPaymentStatus"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	paymentData := <-q.paymentRepositoryQuery.FindPaymentStatusById(ctx, payload.PaymentId)
	if paymentData.Error != nil {
		return nil, paymentData.Error
	}

	if paymentData.Data == nil {
		return nil, errors.BadRequest("payment not found")
	}

	payment, ok := paymentData.Data.(*entity.PaymentHistory)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data payment")
	}

	midtransResp, err := q.midtransRepositoryQuery.GetTransactionStatus(ctx, payment.PaymentId)
	if err != nil {
		return nil, errors.InternalServerError("failed to check payment status")
	}

	result := response.PaymentStatusResp{
		PaymentId:     payment.PaymentId,
		TicketNumber:  payment.Ticket.TicketNumber,
		Bank:          payment.Payment.VaNumbers[0].Bank,
		VaNumber:      payment.Payment.VaNumbers[0].VaNumber,
		Amount:        midtransResp.GrossAmount,
		PaymentStatus: midtransResp.TransactionStatus,
	}

	return &result, nil

}

func (q queryUsecase) FindOrderPayment(origCtx context.Context, payload request.GetOrderPaymentReq) (*response.GetOrderPaymentResp, error) {
	domain := "paymentUsecase-FindPaymentStatus"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	paymentData := <-q.paymentRepositoryQuery.FindPaymentStatusById(ctx, payload.PaymentId)
	if paymentData.Error != nil {
		return nil, paymentData.Error
	}

	if paymentData.Data == nil {
		return nil, errors.BadRequest("payment not found")
	}

	payment, ok := paymentData.Data.(*entity.PaymentHistory)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data payment")
	}

	bankTicketData := <-q.ticketRepositoryQuery.FindBankTicketByTicketNumber(ctx, payment.Ticket.TicketNumber, payment.Ticket.EventId)
	if bankTicketData.Error != nil {
		return nil, bankTicketData.Error
	}

	if bankTicketData.Data == nil {
		return nil, errors.BadRequest("ticket not found")
	}

	bankTicket, ok := bankTicketData.Data.(*ticketEntity.BankTicket)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data ticket")
	}

	if bankTicket.UserId != payload.UserId {
		return nil, errors.BadRequest("user not eligible")
	}

	eventData := <-q.eventRepositoryQuery.FindEventById(ctx, bankTicket.EventId)
	if eventData.Error != nil {
		return nil, bankTicketData.Error
	}

	if eventData.Data == nil {
		return nil, errors.BadRequest("event not found")
	}

	event, ok := eventData.Data.(*eventEntity.Event)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data ticket")
	}

	userCache, err := q.redis.Get(ctx, fmt.Sprintf("%s:%s", constants.RedisKeyGetProfileUser, payload.UserId)).Result()
	if err != nil {
		return nil, errors.InternalServerError("error redis connection")
	}

	var user userDto.UserData
	err = json.Unmarshal([]byte(userCache), &user)
	if err != nil {
		return nil, errors.InternalServerError("cannot parsing redis user")
	}

	var maxWaitTime string
	if payment.Payment.TransactionStatus == constants.Pending && payment.IsValidPayment {
		count := 15
		then := payment.CreatedAt.Local().Add(time.Duration(+count) * time.Minute)
		maxWaitTime = then.Format("2006-01-02 15:04")
	}

	return &response.GetOrderPaymentResp{
		TicketNumber: bankTicket.TicketNumber,
		FullName:     user.Data.FullName,
		TicketType:   bankTicket.TicketType,
		Bank:         payment.Payment.VaNumbers[0].Bank,
		VaNumber:     payment.Payment.VaNumbers[0].VaNumber,
		Amount:       payment.Payment.GrossAmount,
		EventName:    event.Name,
		Country:      event.Country.Name,
		Place:        event.Country.Place,
		OrderTime:    bankTicket.UpdatedAt.Local(),
		MaxWaitTime:  maxWaitTime,
	}, nil
}

func (q queryUsecase) FindPaymentList(origCtx context.Context, payload request.PaymentList) (*response.OrderListResp, error) {
	domain := "paymentUsecase-FindPaymentStatus"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	paymentData := <-q.paymentRepositoryQuery.FindPaymentByUser(ctx, payload)
	if paymentData.Error != nil {
		return nil, paymentData.Error
	}

	if paymentData.Data == nil {
		return nil, errors.BadRequest("payment not found")
	}

	payments, ok := paymentData.Data.(*[]entity.PaymentHistory)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data payment")
	}

	var collectionData = make([]response.OrderList, 0)
	for _, value := range *payments {
		collectionData = append(collectionData, response.OrderList{
			PaymentId:      value.PaymentId,
			TicketNumber:   value.Ticket.TicketNumber,
			TicketType:     value.Ticket.TicketType,
			Amount:         value.Payment.GrossAmount,
			PaymentStatus:  value.Payment.TransactionStatus,
			IsValidPayment: value.IsValidPayment,
		})
	}

	return &response.OrderListResp{
		CollectionData: collectionData,
		MetaData:       helpers.GenerateMetaData(paymentData.Count, int64(len(*payments)), payload.Page, payload.Size),
	}, nil
}
