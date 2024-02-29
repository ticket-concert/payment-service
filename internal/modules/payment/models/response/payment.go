package response

import (
	"payment-service/internal/pkg/constants"
	"time"
)

type PaymentResp struct {
	PaymentId     string `json:"paymentId"`
	PaymentStatus string `json:"paymentStatus"`
}

type PaymentStatusResp struct {
	PaymentId     string `json:"paymentId"`
	TicketNumber  string `json:"ticketNumber"`
	Bank          string `json:"bank"`
	VaNumber      string `json:"vaNumber"`
	Amount        string `json:"amount"`
	PaymentStatus string `json:"paymentStatus"`
}

type GetOrderPaymentResp struct {
	TicketNumber string    `json:"ticketNumber"`
	FullName     string    `json:"fullName"`
	TicketType   string    `json:"ticketType"`
	Bank         string    `json:"bank"`
	VaNumber     string    `json:"vaNumber"`
	Amount       string    `json:"amount"`
	EventName    string    `json:"eventName"`
	Country      string    `json:"country"`
	Place        string    `json:"place"`
	OrderTime    time.Time `json:"orderTime"`
	MaxWaitTime  string    `json:"maxWaitTime"`
}

type OrderList struct {
	PaymentId      string `json:"paymentId"`
	TicketNumber   string `json:"ticketNumber"`
	TicketType     string `json:"ticketType"`
	Amount         string `json:"amount"`
	PaymentStatus  string `json:"paymentStatus"`
	IsValidPayment bool   `json:"isValidPayment"`
}

type OrderListResp struct {
	CollectionData []OrderList
	MetaData       constants.MetaData
}
