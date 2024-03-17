// user_http_handler_test.go

package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"payment-service/internal/modules/payment/handlers"
	"payment-service/internal/modules/payment/models/request"
	"payment-service/internal/modules/payment/models/response"
	"payment-service/internal/pkg/errors"
	mockcert "payment-service/mocks/modules/payment"
	mocklog "payment-service/mocks/pkg/log"
	mockredis "payment-service/mocks/pkg/redis"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type UserHttpHandlerTestSuite struct {
	suite.Suite

	cUC       *mockcert.UsecaseCommand
	cUQ       *mockcert.UsecaseQuery
	cLog      *mocklog.Logger
	validator *validator.Validate
	cRedis    *mockredis.Collections
	handler   *handlers.PaymentHttpHandler
	app       *fiber.App
}

func (suite *UserHttpHandlerTestSuite) SetupTest() {
	suite.cUC = new(mockcert.UsecaseCommand)
	suite.cUQ = new(mockcert.UsecaseQuery)
	suite.cLog = new(mocklog.Logger)
	suite.validator = validator.New()
	suite.cRedis = new(mockredis.Collections)
	suite.handler = &handlers.PaymentHttpHandler{
		PaymentUsecaseCommand: suite.cUC,
		PaymentUsecaseQuery:   suite.cUQ,
		Logger:                suite.cLog,
		Validator:             suite.validator,
	}
	suite.app = fiber.New()
	handlers.InitPaymentHttpHandler(suite.app, suite.cUC, suite.cUQ, suite.cLog, suite.cRedis)
}

func TestUserHttpHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHttpHandlerTestSuite))
}

func (suite *UserHttpHandlerTestSuite) TestCreatePayment() {
	suite.cUC.On("CreatePayment", mock.Anything, mock.Anything).Return(&response.PaymentResp{}, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := request.PaymentReq{
		TicketNumber: "1",
		UserId:       "id",
		EventId:      "id",
		PaymentType:  "bca",
	}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/save", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/save")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.CreatePayment(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestCreatePaymentErrBody() {
	suite.cUC.On("CreatePayment", mock.Anything, mock.Anything).Return(&response.PaymentResp{}, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := request.PaymentReq{
		TicketNumber: "1",
		UserId:       "id",
		EventId:      "id",
		PaymentType:  "bca",
	}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/save", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/save")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	// ctx.Request().SetBody(requestBody)

	err := suite.handler.CreatePayment(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestCreatePaymentErrValidate() {
	suite.cUC.On("CreatePayment", mock.Anything, mock.Anything).Return(&response.PaymentResp{}, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := request.PaymentReq{}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/save", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/save")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.CreatePayment(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestCreatePaymentErr() {
	suite.cUC.On("CreatePayment", mock.Anything, mock.Anything).Return(nil, errors.BadRequest("error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	reqM := request.PaymentReq{
		TicketNumber: "1",
		UserId:       "id",
		EventId:      "id",
		PaymentType:  "bca",
	}
	requestBody, _ := json.Marshal(reqM)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/save", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/save")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.CreatePayment(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestGetPaymentStatus() {
	suite.cUQ.On("FindPaymentStatus", mock.Anything, mock.Anything).Return(&response.PaymentStatusResp{}, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(fiber.MethodGet, "/v1/status?paymentId=1", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().SetRequestURI("/v1/status?paymentId=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetPaymentStatus(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestGetPaymentStatusErrValidate() {
	suite.cUQ.On("FindPaymentStatus", mock.Anything, mock.Anything).Return(&response.PaymentStatusResp{}, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(fiber.MethodGet, "/v1/status", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().SetRequestURI("/v1/status")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetPaymentStatus(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestGetPaymentStatusErr() {
	suite.cUQ.On("FindPaymentStatus", mock.Anything, mock.Anything).Return(nil, errors.BadRequest("error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(fiber.MethodGet, "/v1/status?paymentId=1", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().SetRequestURI("/v1/status?paymentId=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetPaymentStatus(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestCreateTicketOrder() {
	res := "response"
	suite.cUC.On("CreateTicketOrder", mock.Anything, mock.Anything).Return(&res, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(fiber.MethodGet, "/v1/callback/1", nil)
	req.Header.Set("Content-Type", "application/json")

	suite.app.Get("/v1/callback/:paymentId", suite.handler.CreateTicketOrder)
	rs, err := suite.app.Test(req, -1)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), rs)
}

func (suite *UserHttpHandlerTestSuite) TestCreateTicketOrderErrParams() {
	res := "response"
	suite.cUC.On("CreateTicketOrder", mock.Anything, mock.Anything).Return(&res, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(fiber.MethodGet, "/v1/callback", nil)
	req.Header.Set("Content-Type", "application/json")

	suite.app.Get("/v1/callback", suite.handler.CreateTicketOrder)
	_, err := suite.app.Test(req, -1)
	assert.NoError(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestCreateTicketOrderErr() {
	suite.cUC.On("CreateTicketOrder", mock.Anything, mock.Anything).Return(nil, errors.BadRequest("error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(fiber.MethodGet, "/v1/callback/1", nil)
	req.Header.Set("Content-Type", "application/json")

	suite.app.Get("/v1/callback/:paymentId", suite.handler.CreateTicketOrder)
	rs, err := suite.app.Test(req, -1)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), rs)
}

func (suite *UserHttpHandlerTestSuite) TestGetOrderPayment() {
	suite.cUQ.On("FindOrderPayment", mock.Anything, mock.Anything).Return(&response.GetOrderPaymentResp{}, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(fiber.MethodGet, "/v1/order-status?paymentId=1", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/order-status?paymentId=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetOrderPayment(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestGetOrderPaymentErrValidator() {
	suite.cUQ.On("FindOrderPayment", mock.Anything, mock.Anything).Return(&response.GetOrderPaymentResp{}, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(fiber.MethodGet, "/v1/order-status", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/order-status")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetOrderPayment(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestGetOrderPaymentErr() {
	suite.cUQ.On("FindOrderPayment", mock.Anything, mock.Anything).Return(nil, errors.BadRequest("error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(fiber.MethodGet, "/v1/order-status?paymentId=1", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/order-status?paymentId=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetOrderPayment(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestGetPaymentList() {
	suite.cUQ.On("FindPaymentList", mock.Anything, mock.Anything).Return(&response.OrderListResp{}, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(fiber.MethodGet, "/v1/list?page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/list?page=1&size=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetPaymentList(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestGetPaymentListErrParser() {
	suite.cUQ.On("FindPaymentList", mock.Anything, mock.Anything).Return(&response.OrderListResp{}, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(fiber.MethodGet, "/v1/list?page=aa&size=aa", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/list?page=aa&size=aa")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetPaymentList(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestGetPaymentListErrValidate() {
	suite.cUQ.On("FindPaymentList", mock.Anything, mock.Anything).Return(&response.OrderListResp{}, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(fiber.MethodGet, "/v1/list?page=&size=", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/list?page=&size=")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetPaymentList(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *UserHttpHandlerTestSuite) TestGetPaymentListErr() {
	suite.cUQ.On("FindPaymentList", mock.Anything, mock.Anything).Return(nil, errors.BadRequest("error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	req := httptest.NewRequest(fiber.MethodGet, "/v1/list?page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/list?page=1&size=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetPaymentList(ctx)
	assert.Nil(suite.T(), err)
}
