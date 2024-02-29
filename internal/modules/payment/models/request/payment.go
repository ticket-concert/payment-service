package request

type PaymentReq struct {
	TicketNumber string `json:"ticketNumber" validate:"required"`
	UserId       string `json:"userId" validate:"required"`
	EventId      string `json:"eventId" validate:"required"`
	PaymentType  string `json:"paymentType" bson:"paymentType"`
}

type PaymentStatusReq struct {
	PaymentId string `json:"paymentId" validate:"required"`
}

type TicketOrderReq struct {
	PaymentId string `params:"paymentId" validate:"required"`
}

type UpdatePaymentStatusReq struct {
	PaymentId         string `json:"paymentId"`
	TransactionStatus string `json:"transactionStatus"`
}

type GetOrderPaymentReq struct {
	PaymentId string `json:"paymentId" validate:"required"`
	UserId    string `json:"userId" validate:"required"`
}

type PaymentList struct {
	Page   int64  `query:"page" validate:"required"`
	Size   int64  `query:"size" validate:"required"`
	UserId string `query:"userId"`
}
