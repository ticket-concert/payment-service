package request

type TicketReq struct {
	CountryCode string `json:"countryCode" validate:"required"`
}

type UpdateBankTicketReq struct {
	TicketNumber  string `json:"ticketNumber"`
	PaymentStatus string `json:"paymentStatus"`
}
