package entity

import "time"

type Country struct {
	Name  string `json:"name" bson:"name"`
	Code  string `json:"code" bson:"code"`
	City  string `json:"city" bson:"city"`
	Place string `json:"place" bson:"place"`
}

type Ticket struct {
	TicketId       string    `json:"ticketId" bson:"ticketId"`
	EventId        string    `json:"eventId" bson:"eventId"`
	TicketType     string    `json:"ticketType" bson:"ticketType"`
	TicketPrice    int       `json:"ticketPrice" bson:"ticketPrice"`
	TotalQuota     int       `json:"totalQuota" bson:"totalQuota"`
	TotalRemaining int       `json:"totalRemaining" bson:"totalRemaining"`
	ContinentName  string    `json:"continentName" bson:"continentName"`
	ContinentCode  string    `json:"continentCode" bson:"continentCode"`
	Country        Country   `json:"country" bson:"country"`
	Tag            string    `json:"tag" bson:"tag"`
	CreatedAt      time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt" bson:"updatedAt"`
}

type AggregateTotalTicket struct {
	Id                   string `json:"_id" bson:"_id"`
	CountryName          string `json:"countryName" bson:"countryName"`
	TotalAvailableTicket int    `json:"totalAvailableTicket" bson:"totalAvailableTicket"`
}

type BankTicket struct {
	TicketNumber  string    `json:"ticketNumber" bson:"ticketNumber"`
	SeatNumber    int       `json:"seatNumber" bson:"seatNumber"`
	IsUsed        bool      `json:"isUsed" bson:"isUsed"`
	UserId        string    `json:"userId" bson:"userId"`
	QueueId       string    `json:"queueId" bson:"queueId"`
	TicketId      string    `json:"ticketId" bson:"ticketId"`
	EventId       string    `json:"eventId" bson:"eventId"`
	CountryCode   string    `json:"countryCode" bson:"countryCode"`
	Price         int       `json:"price" bson:"price"`
	PaymentStatus string    `json:"paymentStatus" bson:"paymentStatus"`
	TicketType    string    `json:"ticketType" bson:"ticketType"`
	CreatedAt     time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" bson:"updatedAt"`
}
