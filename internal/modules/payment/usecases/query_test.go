// query_test.go
package usecases_test

import (
	"context"
	"testing"
	"time"

	eventEntity "payment-service/internal/modules/event/models/entity"
	"payment-service/internal/modules/payment"
	"payment-service/internal/modules/payment/models/entity"
	"payment-service/internal/modules/payment/models/request"
	"payment-service/internal/modules/payment/models/response"
	uc "payment-service/internal/modules/payment/usecases"
	ticketEntity "payment-service/internal/modules/ticket/models/entity"
	"payment-service/internal/pkg/constants"
	"payment-service/internal/pkg/errors"
	"payment-service/internal/pkg/helpers"
	mockcertEvent "payment-service/mocks/modules/event"
	mockcert "payment-service/mocks/modules/payment"
	mockcertTicket "payment-service/mocks/modules/ticket"
	mocklog "payment-service/mocks/pkg/log"
	mockredis "payment-service/mocks/pkg/redis"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type QueryUsecaseTestSuite struct {
	suite.Suite
	mockPaymentRepositoryQuery  *mockcert.MongodbRepositoryQuery
	mockMidtransRepositoryQuery *mockcert.MidtransRepositoryQuery
	mockTicketRepositoryQuery   *mockcertTicket.MongodbRepositoryQuery
	mockEventRepositoryQuery    *mockcertEvent.MongodbRepositoryQuery
	mockLogger                  *mocklog.Logger
	mockRedis                   *mockredis.Collections
	usecase                     payment.UsecaseQuery
	ctx                         context.Context
}

func (suite *QueryUsecaseTestSuite) SetupTest() {
	suite.mockPaymentRepositoryQuery = &mockcert.MongodbRepositoryQuery{}
	suite.mockMidtransRepositoryQuery = &mockcert.MidtransRepositoryQuery{}
	suite.mockTicketRepositoryQuery = &mockcertTicket.MongodbRepositoryQuery{}
	suite.mockEventRepositoryQuery = &mockcertEvent.MongodbRepositoryQuery{}
	suite.mockRedis = &mockredis.Collections{}
	suite.mockLogger = &mocklog.Logger{}
	suite.ctx = context.Background()
	suite.usecase = uc.NewQueryUsecase(
		suite.mockPaymentRepositoryQuery,
		suite.mockMidtransRepositoryQuery,
		suite.mockTicketRepositoryQuery,
		suite.mockEventRepositoryQuery,
		suite.mockLogger,
		suite.mockRedis,
	)
}
func TestQueryUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(QueryUsecaseTestSuite))
}

func (suite *QueryUsecaseTestSuite) TestFindPaymentStatus() {
	// Arrange
	payload := request.PaymentStatusReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "91238819711",
						Bank:     "bca",
					},
				},
			},
		},
		Error: nil,
	}

	mockMidtransResp := &response.TransactionStatusResponse{
		GrossAmount:       "40",
		TransactionStatus: constants.Pending,
	}
	suite.mockPaymentRepositoryQuery.On("FindPaymentStatusById", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockMidtransResp, nil)
	// Act
	result, err := suite.usecase.FindPaymentStatus(suite.ctx, payload)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindPaymentStatusErr() {
	// Arrange
	payload := request.PaymentStatusReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "91238819711",
						Bank:     "bca",
					},
				},
			},
		},
		Error: errors.BadRequest("error"),
	}

	mockMidtransResp := &response.TransactionStatusResponse{
		GrossAmount:       "40",
		TransactionStatus: constants.Pending,
	}
	suite.mockPaymentRepositoryQuery.On("FindPaymentStatusById", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockMidtransResp, nil)
	// Act
	result, err := suite.usecase.FindPaymentStatus(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindPaymentStatusErrNil() {
	// Arrange
	payload := request.PaymentStatusReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockMidtransResp := &response.TransactionStatusResponse{
		GrossAmount:       "40",
		TransactionStatus: constants.Pending,
	}
	suite.mockPaymentRepositoryQuery.On("FindPaymentStatusById", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockMidtransResp, nil)
	// Act
	result, err := suite.usecase.FindPaymentStatus(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindPaymentStatusErrParse() {
	// Arrange
	payload := request.PaymentStatusReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.Country{
			Name: "name",
		},
		Error: nil,
	}

	mockMidtransResp := &response.TransactionStatusResponse{
		GrossAmount:       "40",
		TransactionStatus: constants.Pending,
	}
	suite.mockPaymentRepositoryQuery.On("FindPaymentStatusById", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(mockMidtransResp, nil)
	// Act
	result, err := suite.usecase.FindPaymentStatus(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindPaymentStatusErrTransaction() {
	// Arrange
	payload := request.PaymentStatusReq{
		PaymentId: "id",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "91238819711",
						Bank:     "bca",
					},
				},
			},
		},
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentStatusById", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockMidtransRepositoryQuery.On("GetTransactionStatus", mock.Anything, mock.Anything).Return(nil, errors.BadRequest("error"))
	// Act
	result, err := suite.usecase.FindPaymentStatus(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindOrderPayment() {
	// Arrange
	payload := request.GetOrderPaymentReq{
		PaymentId: "id",
		UserId:    "userId",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "91238819711",
						Bank:     "bca",
					},
				},
				TransactionStatus: constants.Settlement,
			},
		},
		Error: nil,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
			UpdatedAt:    time.Now(),
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Name:  "name",
				Place: "place",
			},
		},
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentStatusById", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	// Act
	result, err := suite.usecase.FindOrderPayment(suite.ctx, payload)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindOrderPaymentErrPayment() {
	// Arrange
	payload := request.GetOrderPaymentReq{
		PaymentId: "id",
		UserId:    "userId",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "91238819711",
						Bank:     "bca",
					},
				},
				TransactionStatus: constants.Settlement,
			},
		},
		Error: errors.BadRequest("error"),
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
			UpdatedAt:    time.Now(),
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Name:  "name",
				Place: "place",
			},
		},
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentStatusById", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	// Act
	result, err := suite.usecase.FindOrderPayment(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindOrderPaymentErrNil() {
	// Arrange
	payload := request.GetOrderPaymentReq{
		PaymentId: "id",
		UserId:    "userId",
	}

	mockPaymentById := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
			UpdatedAt:    time.Now(),
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Name:  "name",
				Place: "place",
			},
		},
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentStatusById", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	// Act
	result, err := suite.usecase.FindOrderPayment(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindOrderPaymentErrParse() {
	// Arrange
	payload := request.GetOrderPaymentReq{
		PaymentId: "id",
		UserId:    "userId",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.Country{
			Name: "name",
		},
		Error: nil,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
			UpdatedAt:    time.Now(),
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Name:  "name",
				Place: "place",
			},
		},
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentStatusById", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	// Act
	result, err := suite.usecase.FindOrderPayment(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindOrderPaymentErrBank() {
	// Arrange
	payload := request.GetOrderPaymentReq{
		PaymentId: "id",
		UserId:    "userId",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "91238819711",
						Bank:     "bca",
					},
				},
				TransactionStatus: constants.Settlement,
			},
		},
		Error: nil,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
			UpdatedAt:    time.Now(),
		},
		Error: errors.BadRequest("error"),
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Name:  "name",
				Place: "place",
			},
		},
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentStatusById", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	// Act
	result, err := suite.usecase.FindOrderPayment(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindOrderPaymentErrNilBank() {
	// Arrange
	payload := request.GetOrderPaymentReq{
		PaymentId: "id",
		UserId:    "userId",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "91238819711",
						Bank:     "bca",
					},
				},
				TransactionStatus: constants.Settlement,
			},
		},
		Error: nil,
	}

	mockBankByTicketNumber := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Name:  "name",
				Place: "place",
			},
		},
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentStatusById", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	// Act
	result, err := suite.usecase.FindOrderPayment(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindOrderPaymentErrParseBank() {
	// Arrange
	payload := request.GetOrderPaymentReq{
		PaymentId: "id",
		UserId:    "userId",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "91238819711",
						Bank:     "bca",
					},
				},
				TransactionStatus: constants.Settlement,
			},
		},
		Error: nil,
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
			Country: eventEntity.Country{
				Name:  "name",
				Place: "place",
			},
		},
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentStatusById", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	// Act
	result, err := suite.usecase.FindOrderPayment(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindOrderPaymentErrUserId() {
	// Arrange
	payload := request.GetOrderPaymentReq{
		PaymentId: "id",
		UserId:    "userId",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "91238819711",
						Bank:     "bca",
					},
				},
				TransactionStatus: constants.Settlement,
			},
		},
		Error: nil,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "id",
			UpdatedAt:    time.Now(),
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Name:  "name",
				Place: "place",
			},
		},
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentStatusById", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	// Act
	result, err := suite.usecase.FindOrderPayment(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindOrderPaymentErrEvent() {
	// Arrange
	payload := request.GetOrderPaymentReq{
		PaymentId: "id",
		UserId:    "userId",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "91238819711",
						Bank:     "bca",
					},
				},
				TransactionStatus: constants.Settlement,
			},
		},
		Error: nil,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
			UpdatedAt:    time.Now(),
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentStatusById", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	// Act
	result, err := suite.usecase.FindOrderPayment(suite.ctx, payload)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindOrderPaymentErrNilEvent() {
	// Arrange
	payload := request.GetOrderPaymentReq{
		PaymentId: "id",
		UserId:    "userId",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "91238819711",
						Bank:     "bca",
					},
				},
				TransactionStatus: constants.Settlement,
			},
		},
		Error: nil,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
			UpdatedAt:    time.Now(),
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentStatusById", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	// Act
	result, err := suite.usecase.FindOrderPayment(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindOrderPaymentErrParseEvent() {
	// Arrange
	payload := request.GetOrderPaymentReq{
		PaymentId: "id",
		UserId:    "userId",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "91238819711",
						Bank:     "bca",
					},
				},
				TransactionStatus: constants.Settlement,
			},
		},
		Error: nil,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
			UpdatedAt:    time.Now(),
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &entity.Country{
			Name: "name",
		},
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentStatusById", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	// Act
	result, err := suite.usecase.FindOrderPayment(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindOrderPaymentErrRedis() {
	// Arrange
	payload := request.GetOrderPaymentReq{
		PaymentId: "id",
		UserId:    "userId",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "91238819711",
						Bank:     "bca",
					},
				},
				TransactionStatus: constants.Settlement,
			},
		},
		Error: nil,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
			UpdatedAt:    time.Now(),
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Name:  "name",
				Place: "place",
			},
		},
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentStatusById", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("", errors.BadRequest("error")))
	// Act
	result, err := suite.usecase.FindOrderPayment(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindOrderPaymentErrUser() {
	// Arrange
	payload := request.GetOrderPaymentReq{
		PaymentId: "id",
		UserId:    "userId",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "91238819711",
						Bank:     "bca",
					},
				},
				TransactionStatus: constants.Settlement,
			},
		},
		Error: nil,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
			UpdatedAt:    time.Now(),
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Name:  "name",
				Place: "place",
			},
		},
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentStatusById", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult("", nil))
	// Act
	result, err := suite.usecase.FindOrderPayment(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindOrderPaymentErrPending() {
	// Arrange
	payload := request.GetOrderPaymentReq{
		PaymentId: "id",
		UserId:    "userId",
	}

	mockPaymentById := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Ticket: &entity.Ticket{
				TicketNumber: "1",
			},
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "91238819711",
						Bank:     "bca",
					},
				},
				TransactionStatus: constants.Pending,
			},
			IsValidPayment: true,
		},
		Error: nil,
	}

	mockBankByTicketNumber := helpers.Result{
		Data: &ticketEntity.BankTicket{
			TicketNumber: "1",
			UserId:       "userId",
			UpdatedAt:    time.Now(),
		},
		Error: nil,
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Name:  "name",
				Place: "place",
			},
		},
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentStatusById", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentById))
	suite.mockTicketRepositoryQuery.On("FindBankTicketByTicketNumber", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankByTicketNumber))
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRedis.On("Get", suite.ctx, mock.AnythingOfType("string")).Return(redis.NewStringResult(`{"data":{"user_id":"060f949d-36cf-4bb3-8a1b-c4c3a734dafc","full_name":"alif septian","email":"alif_s_nurdianto@telkomsel.co.id","mobileNumber":"+628119621992","role":"user","created_at":"2024-02-27T09:43:17.101Z","updated_at":"2024-02-27T09:43:17.101Z"}}`, nil))
	// Act
	_, err := suite.usecase.FindOrderPayment(suite.ctx, payload)

	// Assert
	assert.NoError(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindPaymentList() {
	// Arrange
	payload := request.PaymentList{
		Page: 1,
		Size: 1,
	}

	mockPaymentByUser := helpers.Result{
		Data: &[]entity.PaymentHistory{
			{
				PaymentId: "id",
				Ticket: &entity.Ticket{
					TicketNumber: "1",
					TicketType:   "Gold",
				},
				Payment: &entity.Payment{
					VaNumbers: []entity.VaNumber{
						{
							VaNumber: "91238819711",
							Bank:     "bca",
						},
					},
					TransactionStatus: constants.Settlement,
					GrossAmount:       "40",
				},
			},
		},
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentByUser", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentByUser))
	// Act
	result, err := suite.usecase.FindPaymentList(suite.ctx, payload)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindPaymentListErr() {
	// Arrange
	payload := request.PaymentList{
		Page: 1,
		Size: 1,
	}

	mockPaymentByUser := helpers.Result{
		Data: &[]entity.PaymentHistory{
			{
				PaymentId: "id",
				Ticket: &entity.Ticket{
					TicketNumber: "1",
					TicketType:   "Gold",
				},
				Payment: &entity.Payment{
					VaNumbers: []entity.VaNumber{
						{
							VaNumber: "91238819711",
							Bank:     "bca",
						},
					},
					TransactionStatus: constants.Settlement,
					GrossAmount:       "40",
				},
			},
		},
		Error: errors.BadRequest("error"),
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentByUser", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentByUser))
	// Act
	result, err := suite.usecase.FindPaymentList(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindPaymentListErrNil() {
	// Arrange
	payload := request.PaymentList{
		Page: 1,
		Size: 1,
	}

	mockPaymentByUser := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentByUser", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentByUser))
	// Act
	result, err := suite.usecase.FindPaymentList(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *QueryUsecaseTestSuite) TestFindPaymentListErrParse() {
	// Arrange
	payload := request.PaymentList{
		Page: 1,
		Size: 1,
	}

	mockPaymentByUser := helpers.Result{
		Data: &entity.PaymentHistory{
			PaymentId: "id",
			Ticket: &entity.Ticket{
				TicketNumber: "1",
				TicketType:   "Gold",
			},
			Payment: &entity.Payment{
				VaNumbers: []entity.VaNumber{
					{
						VaNumber: "91238819711",
						Bank:     "bca",
					},
				},
				TransactionStatus: constants.Settlement,
				GrossAmount:       "40",
			},
		},
		Error: nil,
	}

	suite.mockPaymentRepositoryQuery.On("FindPaymentByUser", mock.Anything, mock.Anything).Return(mockChannel(mockPaymentByUser))
	// Act
	result, err := suite.usecase.FindPaymentList(suite.ctx, payload)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}
