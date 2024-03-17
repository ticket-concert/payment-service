package usecases_test

import (
	"context"
	"payment-service/internal/modules/payment"
	"payment-service/internal/pkg/constants"
	"payment-service/internal/pkg/errors"
	"payment-service/internal/pkg/helpers"
	"testing"

	eventEntity "payment-service/internal/modules/event/models/entity"
	"payment-service/internal/modules/payment/models/entity"
	"payment-service/internal/modules/payment/models/request"
	"payment-service/internal/modules/payment/models/response"
	uc "payment-service/internal/modules/payment/usecases"
	ticketEntity "payment-service/internal/modules/ticket/models/entity"
	userEntity "payment-service/internal/modules/user/models/entity"
	mockcertEvent "payment-service/mocks/modules/event"
	mockcert "payment-service/mocks/modules/payment"
	mockcertTicket "payment-service/mocks/modules/ticket"
	mockcertUser "payment-service/mocks/modules/user"
	mockkafka "payment-service/mocks/pkg/kafka"
	mocklog "payment-service/mocks/pkg/log"
	mockredis "payment-service/mocks/pkg/redis"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CommandUsecaseTestSuite struct {
	suite.Suite
	mockPaymentRepositoryQuery    *mockcert.MongodbRepositoryQuery
	mockPaymentRepositoryCommand  *mockcert.MongodbRepositoryCommand
	mockMidtransRepositoryQuery   *mockcert.MidtransRepositoryQuery
	mockMidtransRepositoryCommand *mockcert.MidtransRepositoryCommand
	mockEventRepositoryQuery      *mockcertEvent.MongodbRepositoryQuery
	mockTicketRepositoryQuery     *mockcertTicket.MongodbRepositoryQuery
	mockTicketRepositoryCommand   *mockcertTicket.MongodbRepositoryCommand
	mockUserRepositoryQuery       *mockcertUser.MongodbRepositoryQuery
	mockLogger                    *mocklog.Logger
	mockKafkaProducer             *mockkafka.Producer
	mockRedis                     *mockredis.Collections
	usecase                       payment.UsecaseCommand
	ctx                           context.Context
}

func (suite *CommandUsecaseTestSuite) SetupTest() {
	suite.mockPaymentRepositoryQuery = &mockcert.MongodbRepositoryQuery{}
	suite.mockPaymentRepositoryCommand = &mockcert.MongodbRepositoryCommand{}
	suite.mockMidtransRepositoryQuery = &mockcert.MidtransRepositoryQuery{}
	suite.mockMidtransRepositoryCommand = &mockcert.MidtransRepositoryCommand{}
	suite.mockEventRepositoryQuery = &mockcertEvent.MongodbRepositoryQuery{}
	suite.mockTicketRepositoryQuery = &mockcertTicket.MongodbRepositoryQuery{}
	suite.mockTicketRepositoryCommand = &mockcertTicket.MongodbRepositoryCommand{}
	suite.mockUserRepositoryQuery = &mockcertUser.MongodbRepositoryQuery{}
	suite.mockLogger = &mocklog.Logger{}
	suite.mockKafkaProducer = &mockkafka.Producer{}
	suite.mockRedis = &mockredis.Collections{}
	suite.ctx = context.Background()
	suite.usecase = uc.NewCommandUsecase(
		suite.mockPaymentRepositoryQuery,
		suite.mockPaymentRepositoryCommand,
		suite.mockMidtransRepositoryQuery,
		suite.mockMidtransRepositoryCommand,
		suite.mockTicketRepositoryQuery,
		suite.mockTicketRepositoryCommand,
		suite.mockEventRepositoryQuery,
		suite.mockUserRepositoryQuery,
		suite.mockKafkaProducer,
		suite.mockLogger,
		suite.mockRedis,
	)
}

func TestCommandUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(CommandUsecaseTestSuite))
}

// Helper function to create a channel
func mockChannel(result helpers.Result) <-chan helpers.Result {
	responseChan := make(chan helpers.Result)

	go func() {
		responseChan <- result
		close(responseChan)
	}()

	return responseChan
}

func (suite *CommandUsecaseTestSuite) TestCreatePayment() {
	payload := request.PaymentReq{
		TicketNumber: "1",
		PaymentType:  "bni",
		UserId:       "userId",
		EventId:      "id",
	}

	mockPaymentByTicketNumber := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockTransferBank := response.BankTransferResponse{
		TransactionStatus: "pending",
		VaNumbers: []response.VaNumber{
			{
				Bank:     "bni",
				VaNumber: "98112081278788",
			},
		},
	}

	mockInsertPayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockPaymentRepositoryQuery.On("FindPaymentByTicketNumber", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentByTicketNumber))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	suite.mockMidtransRepositoryCommand.On("TransferBank", suite.ctx, mock.Anything).Return(&mockTransferBank, nil)
	suite.mockPaymentRepositoryCommand.On("InsertOnePayment", mock.Anything, mock.Anything).Return(mockChannel(mockInsertPayment))
	_, err := suite.usecase.CreatePayment(suite.ctx, payload)
	assert.NoError(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreatePaymentErrBank() {
	payload := request.PaymentReq{
		TicketNumber: "1",
		PaymentType:  "bni",
		UserId:       "userId",
		EventId:      "id",
	}

	mockPaymentByTicketNumber := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockBankByTicketNumber := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}

	mockTransferBank := response.BankTransferResponse{
		TransactionStatus: "pending",
		VaNumbers: []response.VaNumber{
			{
				Bank:     "bni",
				VaNumber: "98112081278788",
			},
		},
	}

	mockInsertPayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockPaymentRepositoryQuery.On("FindPaymentByTicketNumber", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentByTicketNumber))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	suite.mockMidtransRepositoryCommand.On("TransferBank", suite.ctx, mock.Anything).Return(&mockTransferBank, nil)
	suite.mockPaymentRepositoryCommand.On("InsertOnePayment", mock.Anything, mock.Anything).Return(mockChannel(mockInsertPayment))
	_, err := suite.usecase.CreatePayment(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreatePaymentErrNilBank() {
	payload := request.PaymentReq{
		TicketNumber: "1",
		PaymentType:  "bni",
		UserId:       "userId",
		EventId:      "id",
	}

	mockPaymentByTicketNumber := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockBankByTicketNumber := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockTransferBank := response.BankTransferResponse{
		TransactionStatus: "pending",
		VaNumbers: []response.VaNumber{
			{
				Bank:     "bni",
				VaNumber: "98112081278788",
			},
		},
	}

	mockInsertPayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockPaymentRepositoryQuery.On("FindPaymentByTicketNumber", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentByTicketNumber))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	suite.mockMidtransRepositoryCommand.On("TransferBank", suite.ctx, mock.Anything).Return(&mockTransferBank, nil)
	suite.mockPaymentRepositoryCommand.On("InsertOnePayment", mock.Anything, mock.Anything).Return(mockChannel(mockInsertPayment))
	_, err := suite.usecase.CreatePayment(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreatePaymentErrParseBank() {
	payload := request.PaymentReq{
		TicketNumber: "1",
		PaymentType:  "bni",
		UserId:       "userId",
		EventId:      "id",
	}

	mockPaymentByTicketNumber := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockBankByTicketNumber := helpers.Result{
		Data: &entity.Country{
			Name: "name",
		},
		Error: nil,
	}

	mockTransferBank := response.BankTransferResponse{
		TransactionStatus: "pending",
		VaNumbers: []response.VaNumber{
			{
				Bank:     "bni",
				VaNumber: "98112081278788",
			},
		},
	}

	mockInsertPayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockPaymentRepositoryQuery.On("FindPaymentByTicketNumber", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentByTicketNumber))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	suite.mockMidtransRepositoryCommand.On("TransferBank", suite.ctx, mock.Anything).Return(&mockTransferBank, nil)
	suite.mockPaymentRepositoryCommand.On("InsertOnePayment", mock.Anything, mock.Anything).Return(mockChannel(mockInsertPayment))
	_, err := suite.usecase.CreatePayment(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreatePaymentErrUserId() {
	payload := request.PaymentReq{
		TicketNumber: "1",
		PaymentType:  "bni",
		UserId:       "userId",
		EventId:      "id",
	}

	mockPaymentByTicketNumber := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "id",
		},
		Error: nil,
	}

	mockTransferBank := response.BankTransferResponse{
		TransactionStatus: "pending",
		VaNumbers: []response.VaNumber{
			{
				Bank:     "bni",
				VaNumber: "98112081278788",
			},
		},
	}

	mockInsertPayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockPaymentRepositoryQuery.On("FindPaymentByTicketNumber", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentByTicketNumber))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	suite.mockMidtransRepositoryCommand.On("TransferBank", suite.ctx, mock.Anything).Return(&mockTransferBank, nil)
	suite.mockPaymentRepositoryCommand.On("InsertOnePayment", mock.Anything, mock.Anything).Return(mockChannel(mockInsertPayment))
	_, err := suite.usecase.CreatePayment(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreatePaymentErr() {
	payload := request.PaymentReq{
		TicketNumber: "1",
		PaymentType:  "bni",
		UserId:       "userId",
		EventId:      "id",
	}

	mockPaymentByTicketNumber := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}
	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockTransferBank := response.BankTransferResponse{
		TransactionStatus: "pending",
		VaNumbers: []response.VaNumber{
			{
				Bank:     "bni",
				VaNumber: "98112081278788",
			},
		},
	}

	mockInsertPayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockPaymentRepositoryQuery.On("FindPaymentByTicketNumber", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentByTicketNumber))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	suite.mockMidtransRepositoryCommand.On("TransferBank", suite.ctx, mock.Anything).Return(&mockTransferBank, nil)
	suite.mockPaymentRepositoryCommand.On("InsertOnePayment", mock.Anything, mock.Anything).Return(mockChannel(mockInsertPayment))
	_, err := suite.usecase.CreatePayment(suite.ctx, payload)
	assert.NoError(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreatePaymentErrExist() {
	payload := request.PaymentReq{
		TicketNumber: "1",
		PaymentType:  "bni",
		UserId:       "userId",
		EventId:      "id",
	}

	mockPaymentByTicketNumber := helpers.Result{
		Data: &entity.PaymentHistory{
			UserId: "userId",
		},
		Error: nil,
	}
	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockTransferBank := response.BankTransferResponse{
		TransactionStatus: "pending",
		VaNumbers: []response.VaNumber{
			{
				Bank:     "bni",
				VaNumber: "98112081278788",
			},
		},
	}

	mockInsertPayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockPaymentRepositoryQuery.On("FindPaymentByTicketNumber", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentByTicketNumber))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	suite.mockMidtransRepositoryCommand.On("TransferBank", suite.ctx, mock.Anything).Return(&mockTransferBank, nil)
	suite.mockPaymentRepositoryCommand.On("InsertOnePayment", mock.Anything, mock.Anything).Return(mockChannel(mockInsertPayment))
	_, err := suite.usecase.CreatePayment(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreatePaymentErrRedis() {
	payload := request.PaymentReq{
		TicketNumber: "1",
		PaymentType:  "bni",
		UserId:       "userId",
		EventId:      "id",
	}

	mockPaymentByTicketNumber := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockTransferBank := response.BankTransferResponse{
		TransactionStatus: "pending",
		VaNumbers: []response.VaNumber{
			{
				Bank:     "bni",
				VaNumber: "98112081278788",
			},
		},
	}

	mockInsertPayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockPaymentRepositoryQuery.On("FindPaymentByTicketNumber", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentByTicketNumber))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("", errors.BadRequest("error")))
	suite.mockMidtransRepositoryCommand.On("TransferBank", suite.ctx, mock.Anything).Return(&mockTransferBank, nil)
	suite.mockPaymentRepositoryCommand.On("InsertOnePayment", mock.Anything, mock.Anything).Return(mockChannel(mockInsertPayment))
	_, err := suite.usecase.CreatePayment(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreatePaymentErrParseRedis() {
	payload := request.PaymentReq{
		TicketNumber: "1",
		PaymentType:  "bni",
		UserId:       "userId",
		EventId:      "id",
	}

	mockPaymentByTicketNumber := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockTransferBank := response.BankTransferResponse{
		TransactionStatus: "pending",
		VaNumbers: []response.VaNumber{
			{
				Bank:     "bni",
				VaNumber: "98112081278788",
			},
		},
	}

	mockInsertPayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockPaymentRepositoryQuery.On("FindPaymentByTicketNumber", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentByTicketNumber))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("", nil))
	suite.mockMidtransRepositoryCommand.On("TransferBank", suite.ctx, mock.Anything).Return(&mockTransferBank, nil)
	suite.mockPaymentRepositoryCommand.On("InsertOnePayment", mock.Anything, mock.Anything).Return(mockChannel(mockInsertPayment))
	_, err := suite.usecase.CreatePayment(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreatePaymentErrMethod() {
	payload := request.PaymentReq{
		TicketNumber: "1",
		PaymentType:  "mandiri",
		UserId:       "userId",
		EventId:      "id",
	}

	mockPaymentByTicketNumber := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockTransferBank := response.BankTransferResponse{
		TransactionStatus: "pending",
		VaNumbers: []response.VaNumber{
			{
				Bank:     "bni",
				VaNumber: "98112081278788",
			},
		},
	}

	mockInsertPayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockPaymentRepositoryQuery.On("FindPaymentByTicketNumber", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentByTicketNumber))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	suite.mockMidtransRepositoryCommand.On("TransferBank", suite.ctx, mock.Anything).Return(&mockTransferBank, nil)
	suite.mockPaymentRepositoryCommand.On("InsertOnePayment", mock.Anything, mock.Anything).Return(mockChannel(mockInsertPayment))
	_, err := suite.usecase.CreatePayment(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreatePaymentErrTransferBank() {
	payload := request.PaymentReq{
		TicketNumber: "1",
		PaymentType:  "bni",
		UserId:       "userId",
		EventId:      "id",
	}

	mockPaymentByTicketNumber := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockInsertPayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockPaymentRepositoryQuery.On("FindPaymentByTicketNumber", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentByTicketNumber))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	suite.mockMidtransRepositoryCommand.On("TransferBank", suite.ctx, mock.Anything).Return(nil, errors.BadRequest("error"))
	suite.mockPaymentRepositoryCommand.On("InsertOnePayment", mock.Anything, mock.Anything).Return(mockChannel(mockInsertPayment))
	_, err := suite.usecase.CreatePayment(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreatePaymentErrNilTransferBank() {
	payload := request.PaymentReq{
		TicketNumber: "1",
		PaymentType:  "bni",
		UserId:       "userId",
		EventId:      "id",
	}

	mockPaymentByTicketNumber := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockInsertPayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockPaymentRepositoryQuery.On("FindPaymentByTicketNumber", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentByTicketNumber))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	suite.mockMidtransRepositoryCommand.On("TransferBank", suite.ctx, mock.Anything).Return(nil, nil)
	suite.mockPaymentRepositoryCommand.On("InsertOnePayment", mock.Anything, mock.Anything).Return(mockChannel(mockInsertPayment))
	_, err := suite.usecase.CreatePayment(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreatePaymentPermataType() {
	payload := request.PaymentReq{
		TicketNumber: "1",
		PaymentType:  "permata",
		UserId:       "userId",
		EventId:      "id",
	}

	mockPaymentByTicketNumber := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockTransferBank := response.BankTransferResponse{
		TransactionStatus: "pending",
		VaNumbers: []response.VaNumber{
			{
				Bank:     "bni",
				VaNumber: "98112081278788",
			},
		},
		PermataVaNumber: "969639267611",
	}

	mockInsertPayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockPaymentRepositoryQuery.On("FindPaymentByTicketNumber", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentByTicketNumber))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	suite.mockMidtransRepositoryCommand.On("TransferBank", suite.ctx, mock.Anything).Return(&mockTransferBank, nil)
	suite.mockPaymentRepositoryCommand.On("InsertOnePayment", mock.Anything, mock.Anything).Return(mockChannel(mockInsertPayment))
	_, err := suite.usecase.CreatePayment(suite.ctx, payload)
	assert.NoError(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreatePaymentErrInsert() {
	payload := request.PaymentReq{
		TicketNumber: "1",
		PaymentType:  "bni",
		UserId:       "userId",
		EventId:      "id",
	}

	mockPaymentByTicketNumber := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockTransferBank := response.BankTransferResponse{
		TransactionStatus: "pending",
		VaNumbers: []response.VaNumber{
			{
				Bank:     "bni",
				VaNumber: "98112081278788",
			},
		},
	}

	mockInsertPayment := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}

	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockPaymentRepositoryQuery.On("FindPaymentByTicketNumber", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentByTicketNumber))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	suite.mockMidtransRepositoryCommand.On("TransferBank", suite.ctx, mock.Anything).Return(&mockTransferBank, nil)
	suite.mockPaymentRepositoryCommand.On("InsertOnePayment", mock.Anything, mock.Anything).Return(mockChannel(mockInsertPayment))
	_, err := suite.usecase.CreatePayment(suite.ctx, payload)
	assert.NoError(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateTicketOrder() {
	payload := request.TicketOrderReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "92131138822",
						Bank:     "bca",
					},
				},
			},
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
		},
		Error: nil,
	}

	mockOrderByTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockTransactionStatus := &response.TransactionStatusResponse{
		TransactionStatus: constants.Settlement,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
		},
		Error: nil,
	}

	mockUserById := helpers.Result{
		Data: &userEntity.User{
			UserId:       "userId",
			MobileNumber: "081123123123",
		},
		Error: nil,
	}

	mockInsertOrder := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdatePayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentById", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockPaymentRepositoryQuery.On("FindOrderByTicket", mock.Anything, mock.Anything).Return(mockChannel(mockOrderByTicket))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockTransactionStatus, nil)
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockPaymentRepositoryCommand.On("InsertOneOrder", mock.Anything, mock.Anything).Return(mockChannel(mockInsertOrder))
	suite.mockPaymentRepositoryCommand.On("UpdatePaymentStatus", mock.Anything, mock.Anything).Return(mockChannel(mockUpdatePayment))
	suite.mockTicketRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateTicketOrder(suite.ctx, payload)
	assert.NoError(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateTicketOrderErrPayment() {
	payload := request.TicketOrderReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "92131138822",
						Bank:     "bca",
					},
				},
			},
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
		},
		Error: errors.BadRequest("error"),
	}

	mockOrderByTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockTransactionStatus := &response.TransactionStatusResponse{
		TransactionStatus: constants.Settlement,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
		},
		Error: nil,
	}

	mockUserById := helpers.Result{
		Data: &userEntity.User{
			UserId:       "userId",
			MobileNumber: "081123123123",
		},
		Error: nil,
	}

	mockInsertOrder := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdatePayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentById", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockPaymentRepositoryQuery.On("FindOrderByTicket", mock.Anything, mock.Anything).Return(mockChannel(mockOrderByTicket))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockTransactionStatus, nil)
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockPaymentRepositoryCommand.On("InsertOneOrder", mock.Anything, mock.Anything).Return(mockChannel(mockInsertOrder))
	suite.mockPaymentRepositoryCommand.On("UpdatePaymentStatus", mock.Anything, mock.Anything).Return(mockChannel(mockUpdatePayment))
	suite.mockTicketRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateTicketOrder(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateTicketOrderErrNilPayment() {
	payload := request.TicketOrderReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockOrderByTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockTransactionStatus := &response.TransactionStatusResponse{
		TransactionStatus: constants.Settlement,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
		},
		Error: nil,
	}

	mockUserById := helpers.Result{
		Data: &userEntity.User{
			UserId:       "userId",
			MobileNumber: "081123123123",
		},
		Error: nil,
	}

	mockInsertOrder := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdatePayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentById", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockPaymentRepositoryQuery.On("FindOrderByTicket", mock.Anything, mock.Anything).Return(mockChannel(mockOrderByTicket))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockTransactionStatus, nil)
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockPaymentRepositoryCommand.On("InsertOneOrder", mock.Anything, mock.Anything).Return(mockChannel(mockInsertOrder))
	suite.mockPaymentRepositoryCommand.On("UpdatePaymentStatus", mock.Anything, mock.Anything).Return(mockChannel(mockUpdatePayment))
	suite.mockTicketRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateTicketOrder(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateTicketOrderErrParsePayment() {
	payload := request.TicketOrderReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.Country{
			Name: "name",
		},
		Error: nil,
	}

	mockOrderByTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockTransactionStatus := &response.TransactionStatusResponse{
		TransactionStatus: constants.Settlement,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
		},
		Error: nil,
	}

	mockUserById := helpers.Result{
		Data: &userEntity.User{
			UserId:       "userId",
			MobileNumber: "081123123123",
		},
		Error: nil,
	}

	mockInsertOrder := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdatePayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentById", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockPaymentRepositoryQuery.On("FindOrderByTicket", mock.Anything, mock.Anything).Return(mockChannel(mockOrderByTicket))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockTransactionStatus, nil)
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockPaymentRepositoryCommand.On("InsertOneOrder", mock.Anything, mock.Anything).Return(mockChannel(mockInsertOrder))
	suite.mockPaymentRepositoryCommand.On("UpdatePaymentStatus", mock.Anything, mock.Anything).Return(mockChannel(mockUpdatePayment))
	suite.mockTicketRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateTicketOrder(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateTicketOrderErr() {
	payload := request.TicketOrderReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "92131138822",
						Bank:     "bca",
					},
				},
			},
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
		},
		Error: nil,
	}

	mockOrderByTicket := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}

	mockTransactionStatus := &response.TransactionStatusResponse{
		TransactionStatus: constants.Settlement,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
		},
		Error: nil,
	}

	mockUserById := helpers.Result{
		Data: &userEntity.User{
			UserId:       "userId",
			MobileNumber: "081123123123",
		},
		Error: nil,
	}

	mockInsertOrder := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdatePayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentById", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockPaymentRepositoryQuery.On("FindOrderByTicket", mock.Anything, mock.Anything).Return(mockChannel(mockOrderByTicket))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockTransactionStatus, nil)
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockPaymentRepositoryCommand.On("InsertOneOrder", mock.Anything, mock.Anything).Return(mockChannel(mockInsertOrder))
	suite.mockPaymentRepositoryCommand.On("UpdatePaymentStatus", mock.Anything, mock.Anything).Return(mockChannel(mockUpdatePayment))
	suite.mockTicketRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateTicketOrder(suite.ctx, payload)
	assert.NoError(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateTicketOrderErrExist() {
	payload := request.TicketOrderReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "92131138822",
						Bank:     "bca",
					},
				},
			},
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
		},
		Error: nil,
	}

	mockOrderByTicket := helpers.Result{
		Data: &entity.Order{
			OrderId: "id",
		},
		Error: nil,
	}

	mockTransactionStatus := &response.TransactionStatusResponse{
		TransactionStatus: constants.Settlement,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
		},
		Error: nil,
	}

	mockUserById := helpers.Result{
		Data: &userEntity.User{
			UserId:       "userId",
			MobileNumber: "081123123123",
		},
		Error: nil,
	}

	mockInsertOrder := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdatePayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentById", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockPaymentRepositoryQuery.On("FindOrderByTicket", mock.Anything, mock.Anything).Return(mockChannel(mockOrderByTicket))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockTransactionStatus, nil)
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockPaymentRepositoryCommand.On("InsertOneOrder", mock.Anything, mock.Anything).Return(mockChannel(mockInsertOrder))
	suite.mockPaymentRepositoryCommand.On("UpdatePaymentStatus", mock.Anything, mock.Anything).Return(mockChannel(mockUpdatePayment))
	suite.mockTicketRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateTicketOrder(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateTicketOrderErrTransaction() {
	payload := request.TicketOrderReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "92131138822",
						Bank:     "bca",
					},
				},
			},
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
		},
		Error: nil,
	}

	mockOrderByTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
		},
		Error: nil,
	}

	mockUserById := helpers.Result{
		Data: &userEntity.User{
			UserId:       "userId",
			MobileNumber: "081123123123",
		},
		Error: nil,
	}

	mockInsertOrder := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdatePayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentById", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockPaymentRepositoryQuery.On("FindOrderByTicket", mock.Anything, mock.Anything).Return(mockChannel(mockOrderByTicket))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(nil, errors.BadRequest("error"))
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockPaymentRepositoryCommand.On("InsertOneOrder", mock.Anything, mock.Anything).Return(mockChannel(mockInsertOrder))
	suite.mockPaymentRepositoryCommand.On("UpdatePaymentStatus", mock.Anything, mock.Anything).Return(mockChannel(mockUpdatePayment))
	suite.mockTicketRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateTicketOrder(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateTicketOrderErrPending() {
	payload := request.TicketOrderReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "92131138822",
						Bank:     "bca",
					},
				},
			},
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
		},
		Error: nil,
	}

	mockOrderByTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockTransactionStatus := &response.TransactionStatusResponse{
		TransactionStatus: constants.Pending,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
		},
		Error: nil,
	}

	mockUserById := helpers.Result{
		Data: &userEntity.User{
			UserId:       "userId",
			MobileNumber: "081123123123",
		},
		Error: nil,
	}

	mockInsertOrder := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdatePayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentById", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockPaymentRepositoryQuery.On("FindOrderByTicket", mock.Anything, mock.Anything).Return(mockChannel(mockOrderByTicket))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockTransactionStatus, nil)
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockPaymentRepositoryCommand.On("InsertOneOrder", mock.Anything, mock.Anything).Return(mockChannel(mockInsertOrder))
	suite.mockPaymentRepositoryCommand.On("UpdatePaymentStatus", mock.Anything, mock.Anything).Return(mockChannel(mockUpdatePayment))
	suite.mockTicketRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateTicketOrder(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateTicketOrderErrTransactionStatus() {
	payload := request.TicketOrderReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "92131138822",
						Bank:     "bca",
					},
				},
			},
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
		},
		Error: nil,
	}

	mockOrderByTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockTransactionStatus := &response.TransactionStatusResponse{
		TransactionStatus: constants.Expired,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
		},
		Error: nil,
	}

	mockUserById := helpers.Result{
		Data: &userEntity.User{
			UserId:       "userId",
			MobileNumber: "081123123123",
		},
		Error: nil,
	}

	mockInsertOrder := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdatePayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentById", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockPaymentRepositoryQuery.On("FindOrderByTicket", mock.Anything, mock.Anything).Return(mockChannel(mockOrderByTicket))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockTransactionStatus, nil)
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockPaymentRepositoryCommand.On("InsertOneOrder", mock.Anything, mock.Anything).Return(mockChannel(mockInsertOrder))
	suite.mockPaymentRepositoryCommand.On("UpdatePaymentStatus", mock.Anything, mock.Anything).Return(mockChannel(mockUpdatePayment))
	suite.mockTicketRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateTicketOrder(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateTicketOrderErrBank() {
	payload := request.TicketOrderReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "92131138822",
						Bank:     "bca",
					},
				},
			},
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
		},
		Error: nil,
	}

	mockOrderByTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockTransactionStatus := &response.TransactionStatusResponse{
		TransactionStatus: constants.Settlement,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: errors.BadRequest("error"),
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
		},
		Error: nil,
	}

	mockUserById := helpers.Result{
		Data: &userEntity.User{
			UserId:       "userId",
			MobileNumber: "081123123123",
		},
		Error: nil,
	}

	mockInsertOrder := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdatePayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentById", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockPaymentRepositoryQuery.On("FindOrderByTicket", mock.Anything, mock.Anything).Return(mockChannel(mockOrderByTicket))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockTransactionStatus, nil)
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockPaymentRepositoryCommand.On("InsertOneOrder", mock.Anything, mock.Anything).Return(mockChannel(mockInsertOrder))
	suite.mockPaymentRepositoryCommand.On("UpdatePaymentStatus", mock.Anything, mock.Anything).Return(mockChannel(mockUpdatePayment))
	suite.mockTicketRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateTicketOrder(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateTicketOrderErrNilBank() {
	payload := request.TicketOrderReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "92131138822",
						Bank:     "bca",
					},
				},
			},
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
		},
		Error: nil,
	}

	mockOrderByTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockTransactionStatus := &response.TransactionStatusResponse{
		TransactionStatus: constants.Settlement,
	}

	mockBankByTicketNumber := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
		},
		Error: nil,
	}

	mockUserById := helpers.Result{
		Data: &userEntity.User{
			UserId:       "userId",
			MobileNumber: "081123123123",
		},
		Error: nil,
	}

	mockInsertOrder := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdatePayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentById", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockPaymentRepositoryQuery.On("FindOrderByTicket", mock.Anything, mock.Anything).Return(mockChannel(mockOrderByTicket))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockTransactionStatus, nil)
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockPaymentRepositoryCommand.On("InsertOneOrder", mock.Anything, mock.Anything).Return(mockChannel(mockInsertOrder))
	suite.mockPaymentRepositoryCommand.On("UpdatePaymentStatus", mock.Anything, mock.Anything).Return(mockChannel(mockUpdatePayment))
	suite.mockTicketRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateTicketOrder(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateTicketOrderErrParseBank() {
	payload := request.TicketOrderReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "92131138822",
						Bank:     "bca",
					},
				},
			},
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
		},
		Error: nil,
	}

	mockOrderByTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockTransactionStatus := &response.TransactionStatusResponse{
		TransactionStatus: constants.Settlement,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &entity.Country{
			Name: "name",
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
		},
		Error: nil,
	}

	mockUserById := helpers.Result{
		Data: &userEntity.User{
			UserId:       "userId",
			MobileNumber: "081123123123",
		},
		Error: nil,
	}

	mockInsertOrder := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdatePayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentById", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockPaymentRepositoryQuery.On("FindOrderByTicket", mock.Anything, mock.Anything).Return(mockChannel(mockOrderByTicket))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockTransactionStatus, nil)
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockPaymentRepositoryCommand.On("InsertOneOrder", mock.Anything, mock.Anything).Return(mockChannel(mockInsertOrder))
	suite.mockPaymentRepositoryCommand.On("UpdatePaymentStatus", mock.Anything, mock.Anything).Return(mockChannel(mockUpdatePayment))
	suite.mockTicketRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateTicketOrder(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateTicketOrderErrEvent() {
	payload := request.TicketOrderReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "92131138822",
						Bank:     "bca",
					},
				},
			},
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
		},
		Error: nil,
	}

	mockOrderByTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockTransactionStatus := &response.TransactionStatusResponse{
		TransactionStatus: constants.Settlement,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
		},
		Error: errors.BadRequest("error"),
	}

	mockUserById := helpers.Result{
		Data: &userEntity.User{
			UserId:       "userId",
			MobileNumber: "081123123123",
		},
		Error: nil,
	}

	mockInsertOrder := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdatePayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentById", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockPaymentRepositoryQuery.On("FindOrderByTicket", mock.Anything, mock.Anything).Return(mockChannel(mockOrderByTicket))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockTransactionStatus, nil)
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockPaymentRepositoryCommand.On("InsertOneOrder", mock.Anything, mock.Anything).Return(mockChannel(mockInsertOrder))
	suite.mockPaymentRepositoryCommand.On("UpdatePaymentStatus", mock.Anything, mock.Anything).Return(mockChannel(mockUpdatePayment))
	suite.mockTicketRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateTicketOrder(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateTicketOrderErrNilEvent() {
	payload := request.TicketOrderReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "92131138822",
						Bank:     "bca",
					},
				},
			},
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
		},
		Error: nil,
	}

	mockOrderByTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockTransactionStatus := &response.TransactionStatusResponse{
		TransactionStatus: constants.Settlement,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUserById := helpers.Result{
		Data: &userEntity.User{
			UserId:       "userId",
			MobileNumber: "081123123123",
		},
		Error: nil,
	}

	mockInsertOrder := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdatePayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentById", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockPaymentRepositoryQuery.On("FindOrderByTicket", mock.Anything, mock.Anything).Return(mockChannel(mockOrderByTicket))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockTransactionStatus, nil)
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockPaymentRepositoryCommand.On("InsertOneOrder", mock.Anything, mock.Anything).Return(mockChannel(mockInsertOrder))
	suite.mockPaymentRepositoryCommand.On("UpdatePaymentStatus", mock.Anything, mock.Anything).Return(mockChannel(mockUpdatePayment))
	suite.mockTicketRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateTicketOrder(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateTicketOrderErrParseEvent() {
	payload := request.TicketOrderReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "92131138822",
						Bank:     "bca",
					},
				},
			},
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
		},
		Error: nil,
	}

	mockOrderByTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockTransactionStatus := &response.TransactionStatusResponse{
		TransactionStatus: constants.Settlement,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &entity.Country{
			Name: "name",
		},
		Error: nil,
	}

	mockUserById := helpers.Result{
		Data: &userEntity.User{
			UserId:       "userId",
			MobileNumber: "081123123123",
		},
		Error: nil,
	}

	mockInsertOrder := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdatePayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentById", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockPaymentRepositoryQuery.On("FindOrderByTicket", mock.Anything, mock.Anything).Return(mockChannel(mockOrderByTicket))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockTransactionStatus, nil)
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockPaymentRepositoryCommand.On("InsertOneOrder", mock.Anything, mock.Anything).Return(mockChannel(mockInsertOrder))
	suite.mockPaymentRepositoryCommand.On("UpdatePaymentStatus", mock.Anything, mock.Anything).Return(mockChannel(mockUpdatePayment))
	suite.mockTicketRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateTicketOrder(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateTicketOrderErrUser() {
	payload := request.TicketOrderReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "92131138822",
						Bank:     "bca",
					},
				},
			},
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
		},
		Error: nil,
	}

	mockOrderByTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockTransactionStatus := &response.TransactionStatusResponse{
		TransactionStatus: constants.Settlement,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
		},
		Error: nil,
	}

	mockUserById := helpers.Result{
		Data: &userEntity.User{
			UserId:       "userId",
			MobileNumber: "081123123123",
		},
		Error: errors.BadRequest("error"),
	}

	mockInsertOrder := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdatePayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentById", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockPaymentRepositoryQuery.On("FindOrderByTicket", mock.Anything, mock.Anything).Return(mockChannel(mockOrderByTicket))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockTransactionStatus, nil)
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockPaymentRepositoryCommand.On("InsertOneOrder", mock.Anything, mock.Anything).Return(mockChannel(mockInsertOrder))
	suite.mockPaymentRepositoryCommand.On("UpdatePaymentStatus", mock.Anything, mock.Anything).Return(mockChannel(mockUpdatePayment))
	suite.mockTicketRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateTicketOrder(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateTicketOrderErrNilUser() {
	payload := request.TicketOrderReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "92131138822",
						Bank:     "bca",
					},
				},
			},
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
		},
		Error: nil,
	}

	mockOrderByTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockTransactionStatus := &response.TransactionStatusResponse{
		TransactionStatus: constants.Settlement,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
		},
		Error: nil,
	}

	mockUserById := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockInsertOrder := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdatePayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentById", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockPaymentRepositoryQuery.On("FindOrderByTicket", mock.Anything, mock.Anything).Return(mockChannel(mockOrderByTicket))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockTransactionStatus, nil)
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockPaymentRepositoryCommand.On("InsertOneOrder", mock.Anything, mock.Anything).Return(mockChannel(mockInsertOrder))
	suite.mockPaymentRepositoryCommand.On("UpdatePaymentStatus", mock.Anything, mock.Anything).Return(mockChannel(mockUpdatePayment))
	suite.mockTicketRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateTicketOrder(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateTicketOrderErrParseUser() {
	payload := request.TicketOrderReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "92131138822",
						Bank:     "bca",
					},
				},
			},
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
		},
		Error: nil,
	}

	mockOrderByTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockTransactionStatus := &response.TransactionStatusResponse{
		TransactionStatus: constants.Settlement,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
		},
		Error: nil,
	}

	mockUserById := helpers.Result{
		Data: &entity.Country{
			Name: "name",
		},
		Error: nil,
	}

	mockInsertOrder := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdatePayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentById", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockPaymentRepositoryQuery.On("FindOrderByTicket", mock.Anything, mock.Anything).Return(mockChannel(mockOrderByTicket))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockTransactionStatus, nil)
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockPaymentRepositoryCommand.On("InsertOneOrder", mock.Anything, mock.Anything).Return(mockChannel(mockInsertOrder))
	suite.mockPaymentRepositoryCommand.On("UpdatePaymentStatus", mock.Anything, mock.Anything).Return(mockChannel(mockUpdatePayment))
	suite.mockTicketRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateTicketOrder(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateTicketOrderErrInsert() {
	payload := request.TicketOrderReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "92131138822",
						Bank:     "bca",
					},
				},
			},
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
		},
		Error: nil,
	}

	mockOrderByTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockTransactionStatus := &response.TransactionStatusResponse{
		TransactionStatus: constants.Settlement,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
		},
		Error: nil,
	}

	mockUserById := helpers.Result{
		Data: &userEntity.User{
			UserId:       "userId",
			MobileNumber: "081123123123",
		},
		Error: nil,
	}

	mockInsertOrder := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}

	mockUpdatePayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentById", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockPaymentRepositoryQuery.On("FindOrderByTicket", mock.Anything, mock.Anything).Return(mockChannel(mockOrderByTicket))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockTransactionStatus, nil)
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockPaymentRepositoryCommand.On("InsertOneOrder", mock.Anything, mock.Anything).Return(mockChannel(mockInsertOrder))
	suite.mockPaymentRepositoryCommand.On("UpdatePaymentStatus", mock.Anything, mock.Anything).Return(mockChannel(mockUpdatePayment))
	suite.mockTicketRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateTicketOrder(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateTicketOrderErrUpdatePayment() {
	payload := request.TicketOrderReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "92131138822",
						Bank:     "bca",
					},
				},
			},
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
		},
		Error: nil,
	}

	mockOrderByTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockTransactionStatus := &response.TransactionStatusResponse{
		TransactionStatus: constants.Settlement,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
		},
		Error: nil,
	}

	mockUserById := helpers.Result{
		Data: &userEntity.User{
			UserId:       "userId",
			MobileNumber: "081123123123",
		},
		Error: nil,
	}

	mockInsertOrder := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdatePayment := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}

	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentById", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockPaymentRepositoryQuery.On("FindOrderByTicket", mock.Anything, mock.Anything).Return(mockChannel(mockOrderByTicket))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockTransactionStatus, nil)
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockPaymentRepositoryCommand.On("InsertOneOrder", mock.Anything, mock.Anything).Return(mockChannel(mockInsertOrder))
	suite.mockPaymentRepositoryCommand.On("UpdatePaymentStatus", mock.Anything, mock.Anything).Return(mockChannel(mockUpdatePayment))
	suite.mockTicketRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateTicketOrder(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateTicketOrderErrUpdateBank() {
	payload := request.TicketOrderReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "92131138822",
						Bank:     "bca",
					},
				},
			},
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
		},
		Error: nil,
	}

	mockOrderByTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockTransactionStatus := &response.TransactionStatusResponse{
		TransactionStatus: constants.Settlement,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
		},
		Error: nil,
	}

	mockUserById := helpers.Result{
		Data: &userEntity.User{
			UserId:       "userId",
			MobileNumber: "081123123123",
		},
		Error: nil,
	}

	mockInsertOrder := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdatePayment := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentById", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockPaymentRepositoryQuery.On("FindOrderByTicket", mock.Anything, mock.Anything).Return(mockChannel(mockOrderByTicket))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockTransactionStatus, nil)
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockPaymentRepositoryCommand.On("InsertOneOrder", mock.Anything, mock.Anything).Return(mockChannel(mockInsertOrder))
	suite.mockPaymentRepositoryCommand.On("UpdatePaymentStatus", mock.Anything, mock.Anything).Return(mockChannel(mockUpdatePayment))
	suite.mockTicketRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockKafkaProducer.On("Publish", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)

	_, err := suite.usecase.CreateTicketOrder(suite.ctx, payload)
	assert.Error(suite.T(), err)
}
