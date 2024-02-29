package response

type Ticket struct {
	TicketType  string `json:"ticketType"`
	TicketPrice string `json:"ticketPrice"`
	// TotalQuota     string  `json:"totalQuota"`
	// TotalRemaining string  `json:"totalRemaining"`
	ContinentName string `json:"continentName"`
	ContinentCode string `json:"continentCode"`
	CountryName   string `json:"countryName"`
	CountryCode   string `json:"countryCode"`
	IsSold        bool   `json:"isSold"`
}

type TicketResp struct {
	Tickets    []Ticket `json:"tickets"`
	Suggestion []Ticket `json:"suggestion"`
}

type TicketCountry struct {
	CountryName string `json:"countryName"`
	CountryCode string `json:"countryCode"`
	IsSold      bool   `json:"isSold"`
}
