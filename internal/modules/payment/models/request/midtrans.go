package request

type BankTransfer struct {
	Bank     string   `json:"bank"`
	VaNumber string   `json:"va_number"`
	FreeText FreeText `json:"free_text"`
}

type FreeText struct {
	Inquiry []struct {
		ID string `json:"id"`
		EN string `json:"en"`
	} `json:"inquiry"`
	Payment []struct {
		ID string `json:"id"`
		EN string `json:"en"`
	} `json:"payment"`
}

type TransactionDetails struct {
	GrossAmount int    `json:"gross_amount"`
	OrderID     string `json:"order_id"`
}

type CustomerDetails struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

type ItemDetails struct {
	ID       string `json:"id"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
	Name     string `json:"name"`
}

type BankTransferRequest struct {
	PaymentType        string             `json:"payment_type"`
	TransactionDetails TransactionDetails `json:"transaction_details"`
	CustomerDetails    CustomerDetails    `json:"customer_details"`
	ItemDetails        []ItemDetails      `json:"item_details"`
	BankTransfer       BankTransfer       `json:"bank_transfer"`
}

const (
	BankTransferType = "bank_transfer"
	BCA              = "bca"
	BNI              = "bni"
	Permata          = "permata"
)
