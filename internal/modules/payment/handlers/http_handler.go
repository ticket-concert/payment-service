package handlers

import (
	"payment-service/internal/modules/payment"
	"payment-service/internal/modules/payment/models/request"
	"payment-service/internal/pkg/errors"
	"payment-service/internal/pkg/helpers"
	"payment-service/internal/pkg/log"
	"payment-service/internal/pkg/redis"

	middlewares "payment-service/configs/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type PaymentHttpHandler struct {
	PaymentUsecaseCommand payment.UsecaseCommand
	PaymentUsecaseQuery   payment.UsecaseQuery
	Logger                log.Logger
	Validator             *validator.Validate
}

func InitPaymentHttpHandler(app *fiber.App, puc payment.UsecaseCommand, puq payment.UsecaseQuery, log log.Logger, redisClient redis.Collections) {
	handler := &PaymentHttpHandler{
		PaymentUsecaseCommand: puc,
		PaymentUsecaseQuery:   puq,
		Logger:                log,
		Validator:             validator.New(),
	}
	middlewares := middlewares.NewMiddlewares(redisClient)
	route := app.Group("/api/payment")

	route.Post("/v1/save", middlewares.VerifyBearer(), handler.CreatePayment)
	route.Get("/v1/status", middlewares.VerifyBearer(), handler.GetPaymentStatus)
	route.Get("/v1/order-status", middlewares.VerifyBearer(), handler.GetOrderPayment)
	route.Get("/v1/list", middlewares.VerifyBearer(), handler.GetPaymentList)
	route.Get("/v1/callback/:paymentId", middlewares.VerifyBasicAuth(), handler.CreateTicketOrder)
}

func (p PaymentHttpHandler) CreatePayment(c *fiber.Ctx) error {
	req := new(request.PaymentReq)
	if err := c.BodyParser(req); err != nil {
		return helpers.RespError(c, p.Logger, errors.BadRequest("bad request"))
	}

	userId := c.Locals("userId").(string)
	req.UserId = userId

	if err := p.Validator.Struct(req); err != nil {
		return helpers.RespError(c, p.Logger, errors.BadRequest(err.Error()))
	}
	resp, err := p.PaymentUsecaseCommand.CreatePayment(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, p.Logger, err)
	}
	return helpers.RespSuccess(c, p.Logger, resp, "Create payment success")
}

func (p PaymentHttpHandler) GetPaymentStatus(c *fiber.Ctx) error {
	req := new(request.PaymentStatusReq)
	if err := c.QueryParser(req); err != nil {
		return helpers.RespError(c, p.Logger, errors.BadRequest("bad request"))
	}

	if err := p.Validator.Struct(req); err != nil {
		return helpers.RespError(c, p.Logger, errors.BadRequest(err.Error()))
	}
	resp, err := p.PaymentUsecaseQuery.FindPaymentStatus(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, p.Logger, err)
	}
	return helpers.RespSuccess(c, p.Logger, resp, "Get payment status success")
}

func (p PaymentHttpHandler) CreateTicketOrder(c *fiber.Ctx) error {
	req := new(request.TicketOrderReq)
	if err := c.ParamsParser(req); err != nil {
		return helpers.RespError(c, p.Logger, errors.BadRequest("bad request"))
	}

	if err := p.Validator.Struct(req); err != nil {
		return helpers.RespError(c, p.Logger, errors.BadRequest(err.Error()))
	}

	resp, err := p.PaymentUsecaseCommand.CreateTicketOrder(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, p.Logger, err)
	}
	return helpers.RespSuccess(c, p.Logger, resp, "Create payment success")
}

func (p PaymentHttpHandler) GetOrderPayment(c *fiber.Ctx) error {
	req := new(request.GetOrderPaymentReq)
	if err := c.QueryParser(req); err != nil {
		return helpers.RespError(c, p.Logger, errors.BadRequest("bad request"))
	}

	userId := c.Locals("userId").(string)
	req.UserId = userId

	if err := p.Validator.Struct(req); err != nil {
		return helpers.RespError(c, p.Logger, errors.BadRequest(err.Error()))
	}
	resp, err := p.PaymentUsecaseQuery.FindOrderPayment(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, p.Logger, err)
	}
	return helpers.RespSuccess(c, p.Logger, resp, "Get order payment success")
}

func (p PaymentHttpHandler) GetPaymentList(c *fiber.Ctx) error {
	req := new(request.PaymentList)
	if err := c.QueryParser(req); err != nil {
		return helpers.RespError(c, p.Logger, errors.BadRequest("bad request"))
	}

	userId := c.Locals("userId").(string)
	req.UserId = userId

	if err := p.Validator.Struct(req); err != nil {
		return helpers.RespError(c, p.Logger, errors.BadRequest(err.Error()))
	}
	resp, err := p.PaymentUsecaseQuery.FindPaymentList(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, p.Logger, err)
	}
	return helpers.RespPagination(c, p.Logger, resp.CollectionData, resp.MetaData, "Get payment list success")
}
