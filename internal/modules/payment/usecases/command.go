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
	"payment-service/internal/modules/user"
	"regexp"

	"payment-service/internal/pkg/constants"
	"payment-service/internal/pkg/errors"
	"payment-service/internal/pkg/log"
	"payment-service/internal/pkg/redis"
	"time"

	eventEntity "payment-service/internal/modules/event/models/entity"
	ticketEntity "payment-service/internal/modules/ticket/models/entity"
	ticketRequest "payment-service/internal/modules/ticket/models/request"
	userDto "payment-service/internal/modules/user/models/dto"
	userEntity "payment-service/internal/modules/user/models/entity"
	kafkaConfluent "payment-service/internal/pkg/kafka/confluent"

	"github.com/google/uuid"
	"go.elastic.co/apm"
)

type commandUsecase struct {
	paymentRepositoryQuery    payment.MongodbRepositoryQuery
	paymentRepositoryCommand  payment.MongodbRepositoryCommand
	midtransRepositoryQuery   payment.MidtransRepositoryQuery
	midtransRepositoryCommand payment.MidtransRepositoryCommand
	ticketRepositoryQuery     ticket.MongodbRepositoryQuery
	ticketRepositoryCommand   ticket.MongodbRepositoryCommand
	eventRepositoryQuery      event.MongodbRepositoryQuery
	userRepositoryQuery       user.MongodbRepositoryQuery
	kafkaProducer             kafkaConfluent.Producer
	logger                    log.Logger
	redis                     redis.Collections
}

func NewCommandUsecase(
	pmq payment.MongodbRepositoryQuery, pmc payment.MongodbRepositoryCommand, mrq payment.MidtransRepositoryQuery,
	mrc payment.MidtransRepositoryCommand, tmq ticket.MongodbRepositoryQuery, tmc ticket.MongodbRepositoryCommand, emq event.MongodbRepositoryQuery,
	umq user.MongodbRepositoryQuery, kp kafkaConfluent.Producer, log log.Logger, rc redis.Collections) payment.UsecaseCommand {
	return commandUsecase{
		paymentRepositoryQuery:    pmq,
		paymentRepositoryCommand:  pmc,
		midtransRepositoryQuery:   mrq,
		midtransRepositoryCommand: mrc,
		ticketRepositoryQuery:     tmq,
		ticketRepositoryCommand:   tmc,
		eventRepositoryQuery:      emq,
		userRepositoryQuery:       umq,
		kafkaProducer:             kp,
		logger:                    log,
		redis:                     rc,
	}
}

func (c commandUsecase) CreatePayment(origCtx context.Context, payload request.PaymentReq) (*response.PaymentResp, error) {
	domain := "userUsecase-UpdateUser"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	bankTicketData := <-c.ticketRepositoryQuery.FindBankTicketByTicketNumber(ctx, payload.TicketNumber, payload.EventId)
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

	paymentData := <-c.paymentRepositoryQuery.FindPaymentByTicketNumber(ctx, payload.TicketNumber)
	if paymentData.Error != nil {
		return nil, bankTicketData.Error
	}

	if paymentData.Data != nil {
		return nil, errors.BadRequest("ticket on payment")
	}

	userCache, err := c.redis.Get(ctx, fmt.Sprintf("%s:%s", constants.RedisKeyGetProfileUser, payload.UserId)).Result()
	if err != nil {
		return nil, errors.InternalServerError("error redis connection")
	}

	var user userDto.UserData
	err = json.Unmarshal([]byte(userCache), &user)
	if err != nil {
		return nil, errors.InternalServerError("cannot parsing redis user")
	}

	if payload.PaymentType != "bca" && payload.PaymentType != "bni" && payload.PaymentType != "permata" {
		return nil, errors.BadRequest("PaymentType must be bca or bni or permatas")
	}

	mobileNumber := regexp.MustCompile(`^\+628|^628|^08`).ReplaceAllString(user.Data.MobileNumber, "08")
	paymentId := uuid.New().String()

	midTransReq := request.BankTransferRequest{
		PaymentType: request.BankTransferType,
		BankTransfer: request.BankTransfer{
			Bank:     payload.PaymentType,
			VaNumber: mobileNumber,
		},
		TransactionDetails: request.TransactionDetails{
			OrderID:     paymentId,
			GrossAmount: bankTicket.Price,
		},
		CustomerDetails: request.CustomerDetails{
			Email:     user.Data.Email,
			FirstName: user.Data.FullName,
			Phone:     mobileNumber,
		},
		ItemDetails: []request.ItemDetails{
			{
				ID:       bankTicket.TicketNumber,
				Price:    bankTicket.Price,
				Quantity: 1,
				Name:     fmt.Sprintf("Ticket Category %s", bankTicket.TicketType),
			},
		},
	}

	fmt.Println("midTransReq: ", midTransReq)

	midTransResp, err := c.midtransRepositoryCommand.TransferBank(ctx, midTransReq)
	if err != nil {
		return nil, errors.InternalServerError("failed to create payment")
	}

	fmt.Println("midTransResp: ", midTransResp)

	if midTransResp == nil {
		return nil, errors.InternalServerError("unable to process payment")
	}

	var vaNumber []entity.VaNumber

	if payload.PaymentType == request.Permata {
		vaNumber = append(vaNumber, entity.VaNumber{
			Bank:     request.Permata,
			VaNumber: midTransResp.PermataVaNumber,
		})
	} else {
		vaNumber = append(vaNumber, entity.VaNumber{
			Bank:     midTransResp.VaNumbers[0].Bank,
			VaNumber: midTransResp.VaNumbers[0].VaNumber,
		})
	}

	expiryTime := time.Now().Local().Add(15 * time.Minute)
	historyReq := entity.PaymentHistory{
		PaymentId: paymentId,
		UserId:    payload.UserId,
		Ticket: &entity.Ticket{
			TicketNumber: bankTicket.TicketNumber,
			TicketType:   bankTicket.TicketType,
			SeatNumber:   bankTicket.SeatNumber,
			CountryCode:  bankTicket.CountryCode,
			TicketId:     bankTicket.TicketId,
			EventId:      bankTicket.EventId,
		},
		Payment: &entity.Payment{
			StatusCode:        midTransResp.StatusCode,
			TransactionID:     midTransResp.TransactionID,
			GrossAmount:       midTransResp.GrossAmount,
			PaymentType:       midTransResp.PaymentType,
			TransactionStatus: midTransResp.TransactionStatus,
			FraudStatus:       midTransResp.FraudStatus,
			StatusMessage:     midTransResp.StatusMessage,
			MerchantID:        midTransResp.MerchantID,
			PermataVaNumber:   midTransResp.PermataVaNumber,
			VaNumbers:         vaNumber,
		},
		IsValidPayment: true,
		ExpiryTime:     expiryTime,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	historyResp := <-c.paymentRepositoryCommand.InsertOnePayment(ctx, historyReq)
	if historyResp.Error != nil {
		return nil, bankTicketData.Error
	}
	response := response.PaymentResp{
		PaymentId:     historyReq.PaymentId,
		PaymentStatus: midTransResp.TransactionStatus,
	}

	return &response, nil
}

func (c commandUsecase) CreateTicketOrder(origCtx context.Context, payload request.TicketOrderReq) (*string, error) {
	domain := "paymentUsecase-FindPaymentStatus"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	paymentId := payload.PaymentId
	paymentData := <-c.paymentRepositoryQuery.FindPaymentById(ctx, paymentId)
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

	orderData := <-c.paymentRepositoryQuery.FindOrderByTicket(ctx, payment.Ticket.TicketNumber)
	if orderData.Error != nil {
		return nil, paymentData.Error
	}

	if orderData.Data != nil {
		return nil, errors.BadRequest("order ticket paid")
	}

	midtransResp, err := c.midtransRepositoryQuery.GetTransactionStatus(ctx, payment.PaymentId)
	if err != nil {
		return nil, errors.InternalServerError("failed to check payment status")
	}
	if midtransResp.TransactionStatus == constants.Pending {
		return nil, errors.BadRequest("transaction still pending")
	}

	if midtransResp.TransactionStatus != constants.Settlement {
		return nil, errors.BadRequest("transaction not settlement")
	}
	bankTicketData := <-c.ticketRepositoryQuery.FindBankTicketByTicketNumber(ctx, payment.Ticket.TicketNumber, payment.Ticket.EventId)
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

	eventData := <-c.eventRepositoryQuery.FindEventById(ctx, bankTicket.EventId)
	if eventData.Error != nil {
		return nil, eventData.Error
	}

	if eventData.Data == nil {
		return nil, errors.BadRequest("event not found")
	}

	event, ok := eventData.Data.(*eventEntity.Event)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data ticket")
	}

	userData := <-c.userRepositoryQuery.FindOneUserId(ctx, payment.UserId)
	if userData.Error != nil {
		return nil, userData.Error
	}

	if userData.Data == nil {
		return nil, errors.BadRequest("user not found")
	}

	user, ok := userData.Data.(*userEntity.User)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data user")
	}

	orderReq := entity.Order{
		OrderId:       uuid.NewString(),
		PaymentId:     paymentId,
		MobileNumber:  user.MobileNumber,
		VaNumber:      payment.Payment.VaNumbers[0].VaNumber,
		Bank:          payment.Payment.VaNumbers[0].Bank,
		Email:         user.Email,
		FullName:      user.FullName,
		TicketNumber:  payment.Ticket.TicketNumber,
		TicketType:    payment.Ticket.TicketType,
		SeatNumber:    payment.Ticket.SeatNumber,
		EventName:     event.Name,
		Country:       entity.Country(event.Country),
		DateTime:      event.DateTime,
		Description:   event.Description,
		Tag:           event.Tag,
		Amount:        bankTicket.Price,
		PaymentStatus: midtransResp.TransactionStatus,
		OrderTime:     bankTicket.UpdatedAt,
		UserId:        user.UserId,
		QueueId:       bankTicket.QueueId,
		TicketId:      bankTicket.TicketId,
		EventId:       bankTicket.EventId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	orderResp := <-c.paymentRepositoryCommand.InsertOneOrder(ctx, orderReq)
	if orderResp.Error != nil {
		return nil, orderResp.Error
	}

	paymentReq := request.UpdatePaymentStatusReq{
		PaymentId:         paymentId,
		TransactionStatus: midtransResp.TransactionStatus,
	}

	paymentResp := <-c.paymentRepositoryCommand.UpdatePaymentStatus(ctx, paymentReq)
	if paymentResp.Error != nil {
		return nil, paymentResp.Error
	}

	bankTicketReq := ticketRequest.UpdateBankTicketReq{
		TicketNumber:  bankTicket.TicketNumber,
		PaymentStatus: midtransResp.TransactionStatus,
	}
	bankTicketResp := <-c.ticketRepositoryCommand.UpdateBankTicket(ctx, bankTicketReq)
	if bankTicketResp.Error != nil {
		return nil, bankTicketResp.Error
	}

	marshaledKafkaData, _ := json.Marshal(orderReq.OrderId)
	topic := "concert-send-email-pdf"
	c.kafkaProducer.Publish(topic, marshaledKafkaData, nil)
	c.logger.Info(ctx, fmt.Sprintf("Send kafka paymentId : %s", paymentId), fmt.Sprintf("%+v", payload))

	result := "Success create ticket order & payment"
	return &result, nil
}
